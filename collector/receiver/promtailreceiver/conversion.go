package promtailreceiver

import (
	"time"

	"github.com/prometheus/common/model"
	"go.opentelemetry.io/collector/consumer/pdata"
)

func convert(lbs model.LabelSet, ts time.Time, entry string) pdata.Logs {
	ls := pdata.NewLogSlice()
	ls.Resize(1)
	ls.At(0).Body().InitEmpty()

	record := ls.At(0)
	attrs := record.Attributes()
	for k, v := range lbs {
		attrs.InsertString(string(k), string(v))
	}

	record.SetTimestamp(pdata.TimestampUnixNano(ts.UnixNano()))
	record.Body().SetStringVal(entry)

	// package
	out := pdata.NewLogs()
	logs := out.ResourceLogs()
	logs.Resize(1)
	rls := logs.At(0)
	rls.InstrumentationLibraryLogs().Resize(1)
	logSlice := rls.InstrumentationLibraryLogs().At(0).Logs()
	ls.MoveAndAppendTo(logSlice)

	return out
}
