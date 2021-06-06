package checks

import (
	"context"
	"errors"
	"time"

	"github.com/f1shl3gs/manta/pkg/tsdb"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/log"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

const (
	// AlertMetricName is the metric name for synthetic alert timeseries.
	alertMetricName = "ALERTS"
	// AlertForStateMetricName is the metric name for 'for' state of alert.
	alertForStateMetricName = "ALERTS_FOR_STATE"

	// AlertStateLabel is the label name indicating the state of an alert.
	alertStateLabel = "alertstate"
)

// Checker should be implement as a Executor's handler
type Checker struct {
	logger *zap.Logger

	cs     manta.CheckService
	es     manta.EventService
	ts     tsdb.TenantStorage
	engine *promql.Engine
}

func NewChecker(logger *zap.Logger, cs manta.CheckService, es manta.EventService, ts tsdb.TenantStorage) *Checker {
	engOpts := promql.EngineOpts{
		Logger: log.NewZapToGokitLogAdapter(logger),
		// Reg:           prometheus.DefaultRegisterer,
		MaxSamples:    50000000,
		Timeout:       2 * time.Minute,
		LookbackDelta: 5 * time.Minute,
		NoStepSubqueryIntervalFn: func(rangeMillis int64) int64 {
			return time.Minute.Milliseconds()
		},
	}

	return &Checker{
		cs:     cs,
		es:     es,
		ts:     ts,
		engine: promql.NewEngine(engOpts),
		logger: logger,
	}
}

func (checker *Checker) Process(ctx context.Context, task *manta.Task, ts time.Time) error {
	span, ctx := tracing.StartSpanFromContextWithOperationName(ctx, "check")
	defer span.Finish()

	c, err := checker.cs.FindCheckByID(ctx, task.OwnerID)
	if err != nil {
		return err
	}

	span.LogKV(
		"checkID", c.ID.String(),
		"checkName", c.Name,
		"taskID", task.ID.String(),
		"orgID", c.OrgID.String())

	// todo: interpolate
	expr := c.Expr

	qry, err := checker.ts.Queryable(ctx, c.OrgID)
	if err != nil {
		return err
	}

	// reuse time!?
	q, err := checker.engine.NewInstantQuery(qry, expr, ts)
	if err != nil {
		return err
	}
	defer q.Close()

	result := q.Exec(ctx)
	if result.Err != nil {
		return result.Err
	}

	var vector promql.Vector
	switch v := result.Value.(type) {
	case promql.Vector:
		vector = v
	case promql.Scalar:
		vector = promql.Vector{
			promql.Sample{
				Point:  promql.Point(v),
				Metric: labels.Labels{},
			},
		}

	default:
		return errors.New("query result is not a vector or scalar")
	}

	appendable, err := checker.ts.Appendable(ctx, c.OrgID)
	if err != nil {
		return err
	}

	appender := appendable.Appender(ctx)
	defer func() {
		err := appender.Commit()
		if err == nil {
			return
		}

		checker.logger.Warn("commit ALERTS failed",
			zap.Error(err))
	}()

	timestamp := fromTime(time.Now())
	for _, sample := range vector {
		v := sample.V
		for _, condition := range c.Conditions {
			var (
				match     = false
				threshold = condition.Threshold
			)

			switch threshold.Type {
			case manta.GreatThan:
				match = v > threshold.Value
			case manta.Equal:
				match = v == threshold.Value
			case manta.NotEqual:
				match = v != threshold.Value
			case manta.LessThan:
				match = v < threshold.Value
			case manta.Inside:
				match = v > threshold.Max && threshold.Min <= v
			case manta.Outside:
				match = v < threshold.Min && threshold.Max > v
			default:
				checker.logger.Error("Unknown threshold type",
					zap.String("check", c.ID.String()),
					zap.String("type", threshold.Type))
				continue
			}

			if !match {
				continue
			}

			l := make(map[string]string, len(sample.Metric))
			for _, lbl := range sample.Metric {
				l[lbl.Name] = lbl.Value
			}

			lb := labels.NewBuilder(sample.Metric)
			lb.Set(labels.MetricName, alertMetricName)
			lb.Set(labels.AlertName, c.Name)
			lb.Set("check", c.ID.String())

			baseLabels := lb.Labels()

			var vec = promql.Vector{
				valueSample(baseLabels, c, timestamp, v),
				stateSample(baseLabels, c, timestamp),
			}

			// todo: set check's labels
			for _, s := range vec {
				_, err := appender.Append(0, s.Metric, s.T, s.V)
				if err != nil {
					// todo: mark check run failed

					checker.logger.Warn("append ALERTS metrics failed",
						zap.String("check", task.OwnerID.String()),
						zap.Error(err))
				}
			}

			// todo: build annotations

			checker.logger.Debug("Alert",
				zap.String("lb", lb.Labels().String()))
		}
	}

	return nil
}

func fromTime(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/int64(time.Millisecond)
}

func valueSample(lbs labels.Labels, c *manta.Check, ts int64, v float64) promql.Sample {
	lb := labels.NewBuilder(lbs)

	lb.Set(labels.MetricName, alertMetricName)
	lb.Set(labels.AlertName, c.Name)

	return promql.Sample{
		Metric: lb.Labels(),
		Point:  promql.Point{T: ts, V: v},
	}
}

func stateSample(lbs labels.Labels, c *manta.Check, ts int64) promql.Sample {
	lb := labels.NewBuilder(lbs)
	lb.Set(labels.MetricName, alertForStateMetricName)
	lb.Set(labels.AlertName, c.Name)

	return promql.Sample{
		Metric: lb.Labels(),
		Point: promql.Point{
			T: ts,
			V: 1,
		},
	}
}
