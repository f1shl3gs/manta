package web

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/textparse"
	"github.com/prometheus/prometheus/pkg/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/stats"
	v1 "github.com/prometheus/prometheus/web/api/v1"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/log"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/store"
)

type status string

const (
	StatusSuccess status = "success"
	StatusError   status = "error"

	instantQueryPath = "/api/v1/query"
	rangeQueryPath   = "/api/v1/query_range"

	promPrefix          = "/api/v1/query/:orgID/"
	promMetadataPath    = "/api/v1/query/:orgID/api/v1/metadata"
	promLabelNamesPath  = "/api/v1/query/:orgID/api/v1/labels"
	promLabelValuesPath = "/api/v1/query/:orgID/api/v1/label/:name/values"
)

var (
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

	// todo: dummy
	targetRetriever v1.TargetRetriever
}

func NewQueryHandler(logger *zap.Logger, router *Router, tenantStorage store.TenantStorage, tr v1.TargetRetriever) {
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
		logger:          logger,
		Router:          router,
		tenantStorage:   tenantStorage,
		now:             time.Now,
		engine:          engine,
		targetRetriever: tr,
	}

	h.HandlerFunc(http.MethodGet, instantQueryPath, h.handleInstantQuery)
	h.HandlerFunc(http.MethodGet, rangeQueryPath, h.handleRangeQuery)
	h.HandlerFunc(http.MethodGet, promMetadataPath, h.handleMetadata)
	h.HandlerFunc(http.MethodGet, promLabelNamesPath, h.handleLabelNames)
	h.HandlerFunc(http.MethodPost, promLabelNamesPath, h.handleLabelNames)
	h.HandlerFunc(http.MethodGet, promLabelValuesPath, h.handleLabelValues)
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
	typ errorType
	err error
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s: %s", e.typ, e.err)
}

type apiFuncResult struct {
	Data     interface{}      `json:"data,omitempty"`
	Err      *apiError        `json:"err,omitempty"`
	Warnings storage.Warnings `json:"warnings,omitempty"`
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

type metadata struct {
	Type textparse.MetricType `json:"type"`
	Help string               `json:"help"`
	Unit string               `json:"unit"`
}

func (h *QueryHandler) handleMetadata(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		metrics = map[string]map[metadata]struct{}{}
	)

	limit := -1
	if s := r.FormValue("limit"); s != "" {
		var err error
		if limit, err = strconv.Atoi(s); err != nil {
			h.HandleHTTPError(ctx, err, w)
			// apiFuncResult{nil, &apiError{errorBadData, errors.New("limit must be a number")}, nil, nil}
			return
		}
	}

	metric := r.FormValue("metric")
	for _, tt := range h.targetRetriever.TargetsActive() {
		for _, t := range tt {

			if metric == "" {
				for _, mm := range t.MetadataList() {
					m := metadata{Type: mm.Type, Help: mm.Help, Unit: mm.Unit}
					ms, ok := metrics[mm.Metric]

					if !ok {
						ms = map[metadata]struct{}{}
						metrics[mm.Metric] = ms
					}
					ms[m] = struct{}{}
				}
				continue
			}

			if md, ok := t.Metadata(metric); ok {
				m := metadata{Type: md.Type, Help: md.Help, Unit: md.Unit}
				ms, ok := metrics[md.Metric]

				if !ok {
					ms = map[metadata]struct{}{}
					metrics[md.Metric] = ms
				}
				ms[m] = struct{}{}
			}
		}
	}

	// Put the elements from the pseudo-set into a slice for marshaling.
	res := map[string][]metadata{}
	for name, set := range metrics {
		if limit >= 0 && len(res) >= limit {
			break
		}

		s := []metadata{}
		for metadata := range set {
			s = append(s, metadata)
		}
		res[name] = s
	}

	err := encodeResponse(ctx, w, http.StatusOK, apiFuncResult{
		Data: res,
	})
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *QueryHandler) handleLabelNames(w http.ResponseWriter, r *http.Request) {
	var (
		orgID manta.ID

		ctx    = r.Context()
		params = httprouter.ParamsFromContext(ctx)
	)

	err := orgID.DecodeFromString(params.ByName("orgID"))
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	start, err := parseTimeParam(r, "start", minTime)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	end, err := parseTimeParam(r, "end", maxTime)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	matcherSets, err := parseMatchersParam(r.Form["match[]"])
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	queryable, err := h.tenantStorage.Queryable(ctx, orgID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	q, err := queryable.Querier(ctx, timestamp.FromTime(start), timestamp.FromTime(end))
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	defer q.Close()

	var (
		names    = make([]string, 0)
		warnings storage.Warnings
	)

	if len(matcherSets) > 0 {
		hints := &storage.SelectHints{
			Start: timestamp.FromTime(start),
			End:   timestamp.FromTime(end),
			// There is no series function, this token is used for lookups that don't need samples.
			Func: "series",
		}

		labelNamesSet := make(map[string]struct{})
		// Get all series which match matchers
		for _, mset := range matcherSets {
			s := q.Select(false, hints, mset...)
			for s.Next() {
				series := s.At()
				for _, lb := range series.Labels() {
					labelNamesSet[lb.Name] = struct{}{}
				}
			}

			warnings = append(warnings, s.Warnings()...)
			if err := s.Err(); err != nil {
				h.HandleHTTPError(ctx, err, w)
				return
			}
		}

		// Convert the map to an array
		names = make([]string, 0, len(labelNamesSet))
		for key := range labelNamesSet {
			names = append(names, key)
		}

		sort.Strings(names)
	} else {
		names, warnings, err = q.LabelNames()
		if err != nil {
			h.HandleHTTPError(ctx, err, w)
			return
		}
	}

	err = encodeResponse(ctx, w, http.StatusOK, &promAPIResult{
		Data:     names,
		Status:   "success",
		Warnings: warnings,
	})
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func parseMatchersParam(matchers []string) ([][]*labels.Matcher, error) {
	var matcherSets [][]*labels.Matcher
	for _, s := range matchers {
		matchers, err := parser.ParseMetricSelector(s)
		if err != nil {
			return nil, err
		}
		matcherSets = append(matcherSets, matchers)
	}

OUTER:
	for _, ms := range matcherSets {
		for _, lm := range ms {
			if lm != nil && !lm.Matches("") {
				continue OUTER
			}
		}
		return nil, errors.New("match[] must contain at least one non-empty matcher")
	}
	return matcherSets, nil
}

func (h *QueryHandler) handleLabelValues(w http.ResponseWriter, r *http.Request) {
	var (
		orgID  manta.ID
		ctx    = r.Context()
		params = httprouter.ParamsFromContext(ctx)
	)

	name := params.ByName("name")

	err := orgID.DecodeFromString(params.ByName("orgID"))
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if !model.LabelNameRE.MatchString(name) {
		h.HandleHTTPError(ctx, &apiError{errorBadData, errors.Errorf("invalid label name: %q", name)}, w)
		return
	}

	start, err := parseTimeParam(r, "start", minTime)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	end, err := parseTimeParam(r, "end", maxTime)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	matcherSets, err := parseMatchersParam(r.Form["match[]"])
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	queryable, err := h.tenantStorage.Queryable(ctx, orgID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	q, err := queryable.Querier(ctx, timestamp.FromTime(start), timestamp.FromTime(end))
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	defer q.Close()

	var (
		vals     []string
		warnings storage.Warnings
	)
	if len(matcherSets) > 0 {
		var callWarnings storage.Warnings
		labelValuesSet := make(map[string]struct{})
		for _, matchers := range matcherSets {
			vals, callWarnings, err = q.LabelValues(name, matchers...)
			if err != nil {
				// todo: add warnings
				h.HandleHTTPError(ctx, &apiError{errorExec, err}, w)
				return
			}
			warnings = append(warnings, callWarnings...)
			for _, val := range vals {
				labelValuesSet[val] = struct{}{}
			}
		}

		vals = make([]string, 0, len(labelValuesSet))
		for val := range labelValuesSet {
			vals = append(vals, val)
		}
	} else {
		vals, warnings, err = q.LabelValues(name)
		if err != nil {
			// todo: add warnings
			h.HandleHTTPError(ctx, &apiError{errorExec, err}, w)
			return
		}

		if vals == nil {
			vals = []string{}
		}
	}

	sort.Strings(vals)

	err = encodeResponse(ctx, w, http.StatusOK, &promAPIResult{
		Data:     vals,
		Status:   "success",
		Warnings: warnings,
	})
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

type promAPIResult struct {
	Data interface{} `json:"data,omitempty"`
	// Err      *apiError        `json:"err,omitempty"`
	Status   string           `json:"status"`
	Warnings storage.Warnings `json:"warnings,omitempty"`
}

// todo
func decodePromResponse(data interface{}, warnings storage.Warnings) error {
	return nil
}
