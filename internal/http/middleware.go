package http

import (
	"github.com/gin-gonic/gin"
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"net/http"
	"time"
)

type event struct {
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	Msg        string
}

func Logger(logger entity.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}
		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}

		logEvent(logger, &event{
			Path:       path,
			Latency:    time.Since(t),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			Msg:        msg,
		})
	}
}

func logEvent(logger entity.Logger, data *event) {
	log := logger.WithFields(map[string]any{
		"method":    data.Method,
		"path":      data.Path,
		"resp_time": data.Latency,
		"status":    data.StatusCode,
	})

	switch {
	case data.StatusCode >= http.StatusBadRequest && data.StatusCode < http.StatusInternalServerError:
		log.Warn(data.Msg)
	case data.StatusCode >= http.StatusInternalServerError:
		log.Error(data.Msg)
	default:

		log.Info(data.Msg)
	}
}
