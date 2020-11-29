package lokiexporter

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/grafana/loki/pkg/logproto"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"
)

const (
	contentType = "application/x-protobuf"
)

type lokiExporter struct {
	logger *zap.Logger

	client *http.Client
}

func (l *lokiExporter) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (l *lokiExporter) Shutdown(ctx context.Context) error {
	l.client.CloseIdleConnections()
	return nil
}

func (l *lokiExporter) ConsumeLogs(ctx context.Context, ld pdata.Logs) error {
	pr := &logproto.PushRequest{
		Streams: make([]logproto.Stream, 0, ld.LogRecordCount()),
	}

	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rl := rls.At(i)
		if rl.IsNil() {
			continue
		}

		attrs := rl.Resource().Attributes()
		lstrs := make([]string, 0, attrs.Len())
		attrs.ForEach(func(k string, v pdata.AttributeValue) {
			if v.Type() != pdata.AttributeValueSTRING {
				return
			}

			lstrs = append(lstrs, fmt.Sprintf("%s=%q", k, v.StringVal()))
		})
		sort.Strings(lstrs)
		attrsStr := fmt.Sprintf("{%s}", strings.Join(lstrs, ", "))

		ills := rl.InstrumentationLibraryLogs()
		for j := 0; j < ills.Len(); j++ {
			ils := ills.At(j)
			if ils.IsNil() {
				continue
			}

			if !ils.InstrumentationLibrary().IsNil() {
				continue
			}

			logs := ils.Logs()
			stream := logproto.Stream{
				Labels:  attrsStr,
				Entries: make([]logproto.Entry, 0, logs.Len()),
			}

			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)
				if lr.IsNil() {
					continue
				}

				stream.Entries = append(stream.Entries, logproto.Entry{
					Timestamp: time.Unix(0, int64(lr.Timestamp())),
					Line:      lr.Body().StringVal(),
				})
			}

			pr.Streams = append(pr.Streams, stream)
		}
	}

	return l.push(ctx, pr)
}

func (l *lokiExporter) push(ctx context.Context, pr *logproto.PushRequest) error {
	data, err := encode(pr)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(data)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/loki/api/v1/push", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := l.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		l.logger.Warn("server returned unexpected status code",
			zap.Int("code", resp.StatusCode))
	}

	return nil
}

func encode(req *logproto.PushRequest) ([]byte, error) {
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	buf = snappy.Encode(nil, buf)
	return buf, nil
}

func newLogsExporter(cf configmodels.Exporter, logger *zap.Logger, cli *http.Client) (component.LogsExporter, error) {
	return &lokiExporter{logger: logger, client: cli}, nil
}
