package http

import (
	"github.com/gin-gonic/gin"
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"net/http"
)

func (h *Handler) getStatistic(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respError(c, entity.InvalidBodyError.Decorate("ID is required"))
		return
	}

	statistic, err := h.urlService.Get(c, id)
	if err != nil {
		h.respError(c, err)
		return
	}

	c.JSON(http.StatusOK, entity.URLtoStatistic(*statistic))
}
