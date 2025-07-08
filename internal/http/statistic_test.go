package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/marcosfmartins/url_shortener/internal/entity"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetStatistic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		now := time.Now()
		URL := &entity.URL{
			ID:         "1234",
			Hits:       10,
			LastAccess: &now,
		}

		expected := &entity.StatisticDTO{
			Hits:       URL.Hits,
			LastAccess: URL.LastAccess,
		}

		mock := &entity.URLServiceMock{
			GetFn: func(ctx context.Context, ID string) (*entity.URL, error) {
				return URL, nil
			},
		}

		handler := NewHandler(log, mock)

		srv := httptest.NewServer(handler)
		defer srv.Close()

		client := http.DefaultClient
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		res, err := client.Get(fmt.Sprintf("%s/url/%s/statistic", srv.URL, "1234"))
		assert.NoError(t, err)
		defer res.Body.Close()

		jbody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedJson, err := json.Marshal(expected)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, string(expectedJson), string(jbody))
	})

	t.Run("error", func(t *testing.T) {})
}
