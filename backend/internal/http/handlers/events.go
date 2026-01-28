package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/kyleaupton/arrflix/internal/sse"
	"github.com/labstack/echo/v4"
)

type Events struct {
	svc    *service.Services
	broker *sse.Broker
}

func NewEvents(s *service.Services, broker *sse.Broker) *Events {
	return &Events{svc: s, broker: broker}
}

func (h *Events) RegisterProtected(v1 *echo.Group) {
	v1.GET("/events", h.Stream)
}

func (h *Events) Stream(c echo.Context) error {
	ctx := c.Request().Context()

	// Optional filtering: /events?type=a&type=b
	allowed := map[string]bool{}
	for _, t := range c.QueryParams()["type"] {
		if tt := strings.TrimSpace(t); tt != "" {
			allowed[tt] = true
		}
	}
	typeAllowed := func(t string) bool {
		if len(allowed) == 0 {
			return true
		}
		return allowed[t]
	}

	// SSE headers
	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream")
	res.Header().Set(echo.HeaderCacheControl, "no-cache")
	res.Header().Set(echo.HeaderConnection, "keep-alive")
	res.WriteHeader(http.StatusOK)

	flusher, ok := res.Writer.(http.Flusher)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "streaming unsupported"})
	}

	writeEvent := func(eventName, id string, data []byte) error {
		if _, err := fmt.Fprintf(res, "event: %s\n", eventName); err != nil {
			return err
		}
		if id != "" {
			if _, err := fmt.Fprintf(res, "id: %s\n", id); err != nil {
				return err
			}
		}
		if len(data) > 0 {
			// data must be line-prefixed; keep JSON single-line for simplicity
			if _, err := fmt.Fprintf(res, "data: %s\n", string(data)); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(res, "\n"); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	// Ready event
	if typeAllowed("ready") {
		_ = writeEvent("ready", "", []byte(`{"ok":true}`))
	}

	// Initial snapshot (convenience for consumers like downloads page)
	if typeAllowed("download_jobs_snapshot") {
		jobs, err := h.svc.DownloadJobs.List(ctx)
		if err == nil {
			if b, err := json.Marshal(jobs); err == nil {
				_ = writeEvent("download_jobs_snapshot", "", b)
			}
		}
	}

	sub, cancel := h.broker.Subscribe()
	defer cancel()

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-heartbeat.C:
			if typeAllowed("ping") {
				_ = writeEvent("ping", "", []byte(`{"ts":`+fmt.Sprint(time.Now().Unix())+`}`))
			}
		case ev, ok := <-sub:
			if !ok {
				return nil
			}
			if !typeAllowed(ev.Type) {
				continue
			}
			if err := writeEvent(ev.Type, ev.ID, ev.Data); err != nil {
				return nil
			}
		}
	}
}
