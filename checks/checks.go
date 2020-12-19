package checks

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/labelset"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/query"
	"github.com/f1shl3gs/manta/query/promql"
)

// Checker should be implement as a Executor's handler
type Checker struct {
	logger *zap.Logger

	cs manta.CheckService
	es manta.EventService
	ds manta.DatasourceService
}

func NewChecker(logger *zap.Logger, cs manta.CheckService, es manta.EventService, ds manta.DatasourceService) *Checker {
	return &Checker{
		logger: logger,
		cs:     cs,
		es:     es,
		ds:     ds,
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

	querier, err := checker.querierFromDatasource(ctx, c.Datasource)
	if err != nil {
		return err
	}

	results, err := querier.Query(ctx, expr)
	if err != nil {
		return err
	}

	for _, result := range results {
		// conditions are sort by severity
		if len(result.Values) == 0 {
			continue
		}

		current := result.Values[0].Value

		for _, condition := range c.Conditions {
			var (
				match     = false
				threshold = condition.Threshold
			)

			switch threshold.Type {
			case manta.GreatThan:
				match = current > threshold.Value
			case manta.Equal:
				match = current == threshold.Value
			case manta.NotEqual:
				match = current != threshold.Value
			case manta.LessThan:
				match = current < threshold.Value
			case manta.Inside:
				match = current > threshold.Max && threshold.Min <= current
			case manta.Outside:
				match = current < threshold.Min && threshold.Max > current
			default:
				checker.logger.Error("unknown threshold type",
					zap.String("check", c.ID.String()),
					zap.String("type", threshold.Type))
				continue
			}

			if !match {
				continue
			}

			// build the alert
			lbs := make(labelset.LabelSet, 8)
			anns := make(map[string]string, 8)

			for k, v := range result.Labels {
				lbs.Set(k, v)
			}

			for k, v := range task.Annotations {
				anns[k] = v
			}
			anns["check"] = c.ID.String()
			anns["task"] = task.ID.String()
			anns["current"] = strconv.FormatFloat(current, 'f', 2, 64)
			anns["predicate"] = condition.Threshold.Type
			anns["threshold"] = strconv.FormatFloat(condition.Threshold.Value, 'f', 2, 64)

			/*
				now := time.Now()

				_, err := checker.as.UpsertAlert(ctx, &manta.Alert{
					Labels:      lbs,
					Annotations: anns,
					StartsAt:    now,
					EndsAt:      now.Add(time.Minute),
				})

				if err != nil {
					checker.logger.Error("upsert alert failed",
						zap.Error(err))
				}*/
		}
	}

	return nil
}

// todo: add cache
func (checker *Checker) querierFromDatasource(ctx context.Context, id manta.ID) (query.Querier, error) {
	ds, err := checker.ds.FindDatasourceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	switch ds.Type {
	case "prometheus":
		pc := ds.GetPrometheus()
		return promql.New(pc.Url)
	default:
		return nil, fmt.Errorf("unknown datasource type %q", ds.Type)
	}
}
