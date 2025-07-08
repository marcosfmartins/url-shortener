package http

import (
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) postURL(c *gin.Context) {
	dto := entity.URLDTO{}
	err := c.Bind(&dto)
	if err != nil {
		h.respError(c, entity.InvalidBodyError)
		return
	}

	obj, err := h.urlService.Create(c, dto.URL)
	if err != nil {
		h.respError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": obj.ID})
}

func (h *Handler) getRedirectURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respError(c, entity.InvalidBodyError.Decorate("ID is required"))
		return
	}

	URL, err := h.urlService.GetURL(c, id)
	if err != nil {
		c.Status(http.StatusNotFound)
		_, _ = c.Writer.Write([]byte(notFoundContent))
		return
	}

	c.Redirect(http.StatusPermanentRedirect, URL)
}

func (h *Handler) deleteURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respError(c, entity.InvalidBodyError.Decorate("ID is required"))
		return
	}

	err := h.urlService.Delete(c, id)
	if err != nil {
		h.respError(c, entity.NotFoundError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) getURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respError(c, entity.InvalidBodyError.Decorate("ID is required"))
		return
	}

	URL, err := h.urlService.Get(c, id)
	if err != nil {
		h.respError(c, entity.NotFoundError)
		return
	}

	c.JSON(http.StatusOK, URL)
}
