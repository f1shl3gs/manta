package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	scrapePrefix = "/api/v1/scrapes"
	scrapeIDPath = "/api/v1/scrapes/:id"
)

type ScrapeTargetHandler struct {
	*Router

	logger        *zap.Logger
	scrapeService manta.ScraperTargetService
}

func NewScrapeHandler(backend *Backend, logger *zap.Logger) {
	h := &ScrapeTargetHandler{
		Router:        backend.router,
		logger:        logger.With(zap.String("handler", "scrape")),
		scrapeService: backend.ScraperTargetService,
	}

	h.HandlerFunc(http.MethodGet, scrapeIDPath, h.handleGet)
	h.HandlerFunc(http.MethodGet, scrapePrefix, h.handleList)
	h.HandlerFunc(http.MethodPost, scrapePrefix, h.handleCreate)
	h.HandlerFunc(http.MethodDelete, scrapeIDPath, h.handleDelete)
	h.HandlerFunc(http.MethodPatch, scrapeIDPath, h.handlePatch)
}

//
func (h *ScrapeTargetHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := IDFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	s, err := h.scrapeService.FindScraperTargetByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = encodeResponse(ctx, w, http.StatusOK, s)
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

// handleList returns the scrape targets filter by orgID
func (h *ScrapeTargetHandler) handleList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	orgID, err := OrgIdFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	ts, err := h.scrapeService.FindScraperTargets(ctx, manta.ScraperTargetFilter{OrgID: &orgID})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = encodeResponse(ctx, w, http.StatusOK, &ts)
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

// handleCreate create a scrapeTarget
func (h *ScrapeTargetHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var (
		s   = &manta.ScrapeTarget{}
		ctx = r.Context()
	)

	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.scrapeService.CreateScraperTarget(ctx, s)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// handleDelete delete a scrapte target by ID
func (h *ScrapeTargetHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := IDFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.scrapeService.DeleteScraperTarget(ctx, id)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h.HandleHTTPError(ctx, err, w)
}

func (h *ScrapeTargetHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		upd manta.ScraperTargetUpdate
	)

	id, err := IDFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = decodeScrapePatch(r, &upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	scrape, err := h.scrapeService.UpdateScraperTarget(ctx, id, upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, scrape); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeScrapePatch(r *http.Request, upd *manta.ScraperTargetUpdate) error {
	err := json.NewDecoder(r.Body).Decode(upd)
	if err != nil {
		return err
	}

	return nil
}
