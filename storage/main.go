package main

import (
	"github.com/f1shl3gs/manta/log"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/storage/ingest"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	logger, _ := log.New(os.Stdout)
	tsdbOpts := &tsdb.Options{
		MinBlockDuration:  int64(2 * time.Hour / time.Millisecond),
		MaxBlockDuration:  int64(2 * time.Hour / time.Millisecond),
		RetentionDuration: int64(4 * time.Hour / time.Millisecond),
		NoLockfile:        false,
		WALCompression:    true,
	}

	store := ingest.NewMultiTSDB("data", logger, prometheus.DefaultRegisterer, tsdbOpts, labels.FromStrings("foo", "bar"), "tenant_id", nil, false)
	err := store.Open()
	if err != nil {
		panic(err)
	}

	defer store.Close()

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		span, ctx := tracing.StartSpanFromContext(r.Context())
		defer span.Finish()

		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		buf, err := snappy.Decode(nil, compressed)
		if err != nil {
			logger.Warn("decode remote write body failed",
				zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		wr := prompb.WriteRequest{}
		if err := proto.Unmarshal(buf, &wr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tenant := r.Header.Get("X-Scope-OrgID")
		if tenant == "" {
			tenant = "_default"
		}

		app, err := store.TenantAppendable(tenant)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		appender, err := app.Appender(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		numOutOfOrder := 0
		numDuplicates := 0
		numOutOfBounds := 0

		for _, ts := range wr.Timeseries {
			lset := make(labels.Labels, len(ts.Labels))
			for i := range ts.Labels {
				lset[i] = labels.Label{
					Name:  ts.Labels[i].Name,
					Value: ts.Labels[i].Value,
				}
			}

			for _, s := range ts.Samples {
				_, err := appender.Add(lset, s.Timestamp, s.Value)
				switch err {
				case nil:
					continue
				case storage.ErrOutOfOrderSample:
					numOutOfOrder++
					// level.Debug(r.logger).Log("msg", "Out of order sample", "lset", lset.String(), "sample", s.String())
				case storage.ErrDuplicateSampleForTimestamp:
					numDuplicates++
					// level.Debug(r.logger).Log("msg", "Duplicate sample for timestamp", "lset", lset.String(), "sample", s.String())
				case storage.ErrOutOfBounds:
					numOutOfBounds++
					// level.Debug(r.logger).Log("msg", "Out of bounds metric", "lset", lset.String(), "sample", s.String())
				}
			}
		}

		err = appender.Commit()
		if err != nil {
			logger.Warn("commit failed",
				zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	err = http.ListenAndServe(":10080", nil)
	if err != nil {
		panic(err)
	}
}
