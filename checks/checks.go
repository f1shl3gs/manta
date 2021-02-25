package checks

import (
	"context"
	"errors"
	"time"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/log"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/store"
)

const (
	// AlertMetricName is the metric name for synthetic alert timeseries.
	alertMetricName = "ALERTS"
)

// Checker should be implement as a Executor's handler
type Checker struct {
	logger *zap.Logger

	cs     manta.CheckService
	es     manta.EventService
	ts     store.TenantStorage
	engine *promql.Engine
}

func NewChecker(logger *zap.Logger, cs manta.CheckService, es manta.EventService, ts store.TenantStorage) *Checker {
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

func (checker *Checker) Process(ctx context.Context, task *manta.Task) error {
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
	q, err := checker.engine.NewInstantQuery(qry, expr, time.Now())
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
				checker.logger.Error("unknown threshold type",
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

			lb := labels.NewBuilder(sample.Metric).Del(labels.MetricName)

			// todo: set check's labels
			lb.Set("check", c.ID.String())
			lb.Set("org", c.OrgID.String())

			// todo: build annotations

			checker.logger.Info("alert",
				zap.String("lb", lb.Labels().String()))
		}
	}

	return nil
}
