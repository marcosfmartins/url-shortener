package http

import (
	"github.com/marcosfmartins/url_shortener/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine     *gin.Engine
	urlService entity.URLService
}

func NewHandler(logger entity.Logger, urlService entity.URLService) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	handler := &Handler{
		engine:     engine,
		urlService: urlService,
	}

	engine.Use(gin.Recovery())
	engine.Use(Logger(logger))

	handler.registerRoutes()

	return engine
}

func (h *Handler) registerRoutes() {
	h.engine.POST("/url", h.postURL)
	h.engine.GET("/url/:id", h.getURL)
	h.engine.DELETE("/url/:id", h.deleteURL)
	h.engine.GET("/url/:id/statistic", h.getStatistic)

	h.engine.GET("/:id", h.getRedirectURL)
}

func (h *Handler) respError(c *gin.Context, err error) {
	rawErro, ok := err.(*entity.Error)
	if ok {
		c.JSON(rawErro.Code, gin.H{
			"message": err.Error(),
			"code":    rawErro.Code,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
			"detail":  err.Error(),
			"code":    http.StatusInternalServerError,
		})
	}

	c.Abort()
}
