package web

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/cache"
)

const (
	DatasourcePrefix    = "/api/v1/datasources"
	DatasourceIDPath    = "/api/v1/datasources/:id"
	DatasourceProxyPath = "/api/v1/datasources/:id/proxy/*path"
)

// todo: datasource cache

type DatasourceHandler struct {
	*Router

	logger            *zap.Logger
	cached            *cache.LRU
	singleflight      singleflight.Group
	datasourceService manta.DatasourceService
}

func (h *DatasourceHandler) proxy(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		params = httprouter.ParamsFromContext(ctx)
		key    = params.ByName("id")
		err    error
	)

	val := h.cached.Get(key)
	if val == nil {
		val, err, _ = h.singleflight.Do(key, func() (interface{}, error) {
			var (
				id  manta.ID
				err error
			)

			if err = id.DecodeFromString(key); err != nil {
				return nil, err
			}

			ds, err := h.datasourceService.FindDatasourceByID(ctx, id)
			if err != nil {
				return nil, err
			}

			proxy, err := NewDatasourceProxy(ds)
			if err != nil {
				return nil, err
			}

			return proxy, nil
		})

		if err != nil {
			h.HandleHTTPError(ctx, err, w)
			return
		}
	}

	if val == nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	proxy := val.(*httputil.ReverseProxy)
	proxy.ServeHTTP(w, r)
}

func (h *DatasourceHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromRequestPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	ds, err := h.datasourceService.FindDatasourceByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, ds); err != nil {
		logEncodingError(h.logger, r, err)
		return
	}
}

func (h *DatasourceHandler) handleList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	list, err := h.datasourceService.FindDatasources(ctx, manta.DatasourceFilter{
		OrgID: &orgID,
	})

	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, &list); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeDatasource(r *http.Request) (*manta.Datasource, error) {
	ds := &manta.Datasource{}
	err := json.NewDecoder(r.Body).Decode(ds)
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Op:   "decode datasource",
			Err:  err,
		}
	}

	return ds, nil
}
func (h *DatasourceHandler) createDatasource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ds, err := decodeDatasource(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.datasourceService.CreateDatasource(ctx, ds)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DatasourceService(logger *zap.Logger, router *Router, ds manta.DatasourceService) {
	h := &DatasourceHandler{
		Router:            router,
		datasourceService: ds,
		cached: cache.NewLRUWithOptions(128, &cache.Options{
			TTL:             10 * time.Second,
			InitialCapacity: 16,
		}),
	}

	router.HandlerFunc(http.MethodPut, DatasourcePrefix, h.createDatasource)
	router.HandlerFunc(http.MethodGet, DatasourceProxyPath, h.proxy)
	router.HandlerFunc(http.MethodGet, DatasourceIDPath, h.get)
	router.HandlerFunc(http.MethodGet, DatasourcePrefix, h.handleList)
}

func NewDatasourceProxy(ds *manta.Datasource) (*httputil.ReverseProxy, error) {
	var (
		proxy *httputil.ReverseProxy
	)

	switch ds.Type {
	case "loki":
		cf := ds.GetLoki()
		target, err := url.Parse(cf.Url)
		if err != nil {
			return nil, err
		}

		proxy = httputil.NewSingleHostReverseProxy(target)
		targetQuery := target.RawQuery
		prefix := DatasourcePrefix + "/" + ds.ID.String() + "/proxy"

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}

			req.URL.Path = strings.Replace(req.URL.Path, prefix, "", 1)

			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		}

		return proxy, nil

	default:
		return nil, errors.Errorf("unknown datasource type %s", ds.Type)
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
