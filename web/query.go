package web

import (
	"context"
	"github.com/f1shl3gs/manta/log"
	"github.com/prometheus/client_golang/prometheus"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/stats"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/store"
)

type status string

const (
	StatusSuccess status = "success"
	StatusError   status = "error"
)

var (
	instantQueryPath = "/api/v1/query"
	rangeQueryPath   = "/api/v1/query_range"

	minTime          = time.Unix(math.MinInt64/1000+62135596801, 0).UTC()
	maxTime          = time.Unix(math.MaxInt64/1000-62135596801, 999999999).UTC()
	minTimeFormatted = minTime.Format(time.RFC3339Nano)
	maxTimeFormatted = maxTime.Format(time.RFC3339Nano)
)

type QueryHandler struct {
	*Router

	now func() time.Time

	logger        *zap.Logger
	engine        *promql.Engine
	tenantStorage store.TenantStorage
}

func NewQueryHandler(logger *zap.Logger, router *Router, tenantStorage store.TenantStorage) {
	engOpts := promql.EngineOpts{
		Logger:        log.NewZapToGokitLogAdapter(logger),
		Reg:           prometheus.DefaultRegisterer,
		MaxSamples:    50000000,
		Timeout:       2 * time.Minute,
		LookbackDelta: 5 * time.Minute,
		NoStepSubqueryIntervalFn: func(rangeMillis int64) int64 {
			return time.Minute.Milliseconds()
		},
	}
	engine := promql.NewEngine(engOpts)

	h := &QueryHandler{
		logger:        logger,
		Router:        router,
		tenantStorage: tenantStorage,
		now:           time.Now,
		engine:        engine,
	}

	h.HandlerFunc(http.MethodGet, instantQueryPath, h.handleInstantQuery)
	h.HandlerFunc(http.MethodGet, rangeQueryPath, h.handleRangeQuery)
}

func parseTime(s string) (time.Time, error) {
	if t, err := strconv.ParseFloat(s, 64); err == nil {
		s, ns := math.Modf(t)
		ns = math.Round(ns*1000) / 1000
		return time.Unix(int64(s), int64(ns*float64(time.Second))).UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}

	// Stdlib's time parser can only handle 4 digit years. As a workaround until
	// that is fixed we want to at least support our own boundary times.
	// Context: https://github.com/prometheus/client_golang/issues/614
	// Upstream issue: https://github.com/golang/go/issues/20555
	switch s {
	case minTimeFormatted:
		return minTime, nil
	case maxTimeFormatted:
		return maxTime, nil
	}
	return time.Time{}, errors.Errorf("cannot parse %q to a valid timestamp", s)
}

func parseTimeParam(r *http.Request, param string, defaultVal time.Time) (time.Time, error) {
	val := r.FormValue(param)
	if val == "" {
		return defaultVal, nil
	}

	t, err := parseTime(val)
	if err != nil {
		return time.Time{}, err
	}

	return t, err
}

func (h *QueryHandler) handleInstantQuery(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	ts, err := parseTimeParam(r, "time", h.now())

	queryable, err := h.tenantStorage.Queryable(ctx, orgID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	qry, err := h.engine.NewInstantQuery(queryable, r.FormValue("query"), ts)
	if err != nil {
		h.HandleHTTPError(ctx, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "bad_data: invalid parameter \"query\"",
			Op:   "create instant query",
			Err:  err,
		}, w)
		return
	}

	defer qry.Close()

	h.encodeQueryResult(ctx, w, r, qry)
}

func (h *QueryHandler) handleRangeQuery(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracing.StartSpanFromContext(r.Context())
	defer span.Finish()

	start, err := parseTime(r.FormValue("start"))
	if err != nil {
		h.handleInvalidParam(ctx, w, r, err)
		return
	}

	end, err := parseTime(r.FormValue("end"))
	if err != nil {
		h.handleInvalidParam(ctx, w, r, err)
		return
	}

	if end.Before(start) {
		h.handleInvalidParam(ctx, w, r, errors.New("end timestamp must not be before start time"))
		return
	}

	step, err := parseDuration(r.FormValue("step"))
	if err != nil {
		h.handleInvalidParam(ctx, w, r, err)
		return
	}

	if step <= 0 {
		h.handleInvalidParam(ctx, w, r, errors.New("zero or negative query resolution step widths are not accepted. Try a positive integer"))
		return
	}

	// For safety, limit the number of returned points per timeseries.
	// This is sufficient for 60s resolution for a week or 1h resolution for a year
	if end.Sub(start)/step > 11000 {
		// todo
	}

	if timeout := r.FormValue("timeout"); timeout != "" {
		var cancel context.CancelFunc

		d, err := parseDuration(timeout)
		if err != nil {
			h.handleInvalidParam(ctx, w, r, errors.Wrap(err, "parse timeout failed"))
			return
		}

		ctx, cancel = context.WithTimeout(ctx, d)
		defer cancel()
	}

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	queryable, err := h.tenantStorage.Queryable(ctx, orgID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
	}

	qs := r.FormValue("query")
	span.LogKV("query", qs)

	qry, err := h.engine.NewRangeQuery(queryable, qs, start, end, step)
	if err != nil {
		h.handleInvalidParam(ctx, w, r, errors.Wrap(err, "invalid query"))
		return
	}

	h.encodeQueryResult(ctx, w, r, qry)
}

func (h *QueryHandler) handleInvalidParam(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	h.HandleHTTPError(ctx, &manta.Error{
		Code: manta.EInvalid,
		Msg:  "invalid param",
		Err:  err,
	}, w)
}

func parseDuration(s string) (time.Duration, error) {
	if d, err := strconv.ParseFloat(s, 64); err == nil {
		ts := d * float64(time.Second)
		if ts > float64(math.MaxInt64) || ts < float64(math.MinInt64) {
			return 0, errors.Errorf("cannot parse %q to a valid duration. It overflows int64", s)
		}
		return time.Duration(ts), nil
	}
	if d, err := model.ParseDuration(s); err == nil {
		return time.Duration(d), nil
	}
	return 0, errors.Errorf("cannot parse %q to a valid duration", s)
}

type queryData struct {
	ResultType parser.ValueType  `json:"resultType"`
	Result     parser.Value      `json:"result"`
	Stats      *stats.QueryStats `json:"stats,omitempty"`
}

type errorType string

const (
	errorNone        errorType = ""
	errorTimeout     errorType = "timeout"
	errorCanceled    errorType = "canceled"
	errorExec        errorType = "execution"
	errorBadData     errorType = "bad_data"
	errorInternal    errorType = "internal"
	errorUnavailable errorType = "unavailable"
	errorNotFound    errorType = "not_found"
)

type apiError struct {
	Typ errorType `json:"typ"`
	Err error     `json:"err"`
}

type apiFuncResult struct {
	Data      interface{}      `json:"data,omitempty"`
	Err       *apiError        `json:"err,omitempty"`
	Warnings  storage.Warnings `json:"warnings,omitempty"`
	finalizer func()
}

func returnAPIError(err error) *apiError {
	if err == nil {
		return nil
	}

	switch errors.Cause(err).(type) {
	case promql.ErrQueryCanceled:
		return &apiError{errorCanceled, err}
	case promql.ErrQueryTimeout:
		return &apiError{errorTimeout, err}
	case promql.ErrStorage:
		return &apiError{errorInternal, err}
	}

	return &apiError{errorExec, err}
}

type ErrorType string

const (
	ErrorNone     ErrorType = ""
	ErrorTimeout  ErrorType = "timeout"
	ErrorCanceled ErrorType = "canceled"
	ErrorExec     ErrorType = "execution"
	ErrorBadData  ErrorType = "bad_data"
	ErrorInternal ErrorType = "internal"
)

type response struct {
	Status    status      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	ErrorType ErrorType   `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
	Warnings  []string    `json:"warnings,omitempty"`
}

func (h *QueryHandler) encodeQueryResult(ctx context.Context, w http.ResponseWriter, r *http.Request, qry promql.Query) {
	res := qry.Exec(ctx)
	defer qry.Close()

	if res.Err != nil {
		resp := response{
			Status: StatusError,
			Error:  res.Err.Error(),
		}

		switch res.Err.(type) {
		case promql.ErrQueryCanceled:
			resp.ErrorType = ErrorCanceled
		case promql.ErrQueryTimeout:
			resp.ErrorType = ErrorTimeout
		case promql.ErrStorage:
			resp.ErrorType = ErrorInternal
		default:
			resp.ErrorType = ErrorExec
		}

		err := encodeResponse(ctx, w, http.StatusOK, resp)
		if err != nil {
			logEncodingError(h.logger, r, err)
		}

		return
	}

	var qs *stats.QueryStats
	if r.FormValue("stats") != "" {
		qs = stats.NewQueryStats(qry.Stats())
	}

	resp := response{
		Data: &queryData{
			ResultType: res.Value.Type(),
			Result:     res.Value,
			Stats:      qs,
		},
		Status: StatusSuccess,
	}

	for _, warn := range res.Warnings {
		resp.Warnings = append(resp.Warnings, warn.Error())
	}

	err := encodeResponse(ctx, w, http.StatusOK, resp)
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}
