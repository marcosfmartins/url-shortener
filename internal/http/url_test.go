package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/marcosfmartins/url_shortener/internal/entity"
	"github.com/marcosfmartins/url_shortener/pkg/logger"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	log = logger.NewZerologAdapter()
)

func TestPostURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &entity.URLServiceMock{
			CreateFn: func(ctx context.Context, url string) (*entity.URL, error) {
				return &entity.URL{ID: "123", Original: url}, nil
			},
		}

		handler := NewHandler(log, mock)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		data := entity.URLDTO{URL: "https://example.com"}

		jsonData, err := json.Marshal(data)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Post(srv.URL+"/url", "application/json", bytes.NewReader(jsonData))
		assert.NoError(t, err)

		bdata, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		defer res.Body.Close()

		var content map[string]string
		err = json.Unmarshal(bdata, &content)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "123", content["id"])
	})

	t.Run("invalid body", func(t *testing.T) {
		handler := NewHandler(log, nil)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		res, err := http.DefaultClient.Post(srv.URL+"/url", "application/json", nil)
		assert.NoError(t, err)

		bdata, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.JSONEq(t, `{"code":400,"message":"invalid request body"}`, string(bdata))
	})

	t.Run("service error", func(t *testing.T) {
		mock := &entity.URLServiceMock{
			CreateFn: func(ctx context.Context, url string) (*entity.URL, error) {
				return nil, assert.AnError
			},
		}

		handler := NewHandler(log, mock)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		data := entity.URLDTO{URL: "https://example.com"}

		jsonData, err := json.Marshal(data)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Post(srv.URL+"/url", "application/json", bytes.NewReader(jsonData))
		assert.NoError(t, err)

		bdata, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.JSONEq(t, `{"code":500, "detail":"assert.AnError general error for testing", "message":"internal server error"}`, string(bdata))
	})
}

func TestGetRedirectURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		redirectURL := "https://example.com"
		mock := &entity.URLServiceMock{
			GetURLFn: func(ctx context.Context, ID string) (string, error) {
				return redirectURL, nil
			},
		}

		handler := NewHandler(log, mock)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		client := http.DefaultClient
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		res, err := client.Get(fmt.Sprintf("%s/%s", srv.URL, "1234"))
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusPermanentRedirect, res.StatusCode)
		assert.Equal(t, redirectURL, res.Header.Get("Location"))
	})

	t.Run("service error", func(t *testing.T) {
		mock := &entity.URLServiceMock{
			GetURLFn: func(ctx context.Context, ID string) (string, error) {
				return "", assert.AnError
			},
		}

		handler := NewHandler(log, mock)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, "1234"))
		assert.NoError(t, err)

		bdata, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Contains(t, string(bdata), `<title>Página não encontrada</title>`)
	})
}

func TestDeleteURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &entity.URLServiceMock{
			DeleteFn: func(ctx context.Context, ID string) error {
				return nil
			},
		}

		handler := NewHandler(log, mock)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/url/%s", srv.URL, "1234"), nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		mock := &entity.URLServiceMock{
			DeleteFn: func(ctx context.Context, ID string) error {
				return assert.AnError
			},
		}

		handler := NewHandler(log, mock)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/url/%s", srv.URL, "1234"), nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer res.Body.Close()

		bdata, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.JSONEq(t, `{"code":404,"message":"not found"}`, string(bdata))
	})
}

func TestGetUrl(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		obj := &entity.URL{ID: "1234", Original: "https://example.com"}

		mock := &entity.URLServiceMock{
			GetFn: func(ctx context.Context, ID string) (*entity.URL, error) {
				return obj, nil
			},
		}

		handler := NewHandler(log, mock)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/url/%s", srv.URL, "1234"))
		assert.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		bdata, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		expected, err := json.Marshal(obj)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expected), string(bdata))
	})
}
