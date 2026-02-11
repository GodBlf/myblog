package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getCounterValue(t *testing.T, method, path, status string) float64 {
	t.Helper()
	metric := &dto.Metric{}
	err := requestTotal.WithLabelValues(method, path, status).Write(metric)
	require.NoError(t, err)
	return metric.GetCounter().GetValue()
}

func TestMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("normal path increases request counter", func(t *testing.T) {
		router := gin.New()
		router.Use(Metric())
		router.GET("/ping", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "pong")
		})

		before := getCounterValue(t, http.MethodGet, "/ping", "200")

		writer := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/ping", nil)
		router.ServeHTTP(writer, request)

		require.Equal(t, http.StatusOK, writer.Code)
		after := getCounterValue(t, http.MethodGet, "/ping", "200")
		assert.Equal(t, before+1, after)
	})

	t.Run("metrics path is skipped", func(t *testing.T) {
		router := gin.New()
		router.Use(Metric())
		router.GET("/metrics", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		before := getCounterValue(t, http.MethodGet, "/metrics", "200")

		writer := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		router.ServeHTTP(writer, request)

		require.Equal(t, http.StatusOK, writer.Code)
		after := getCounterValue(t, http.MethodGet, "/metrics", "200")
		assert.Equal(t, before, after)
	})

	t.Run("unknown path uses unknown label", func(t *testing.T) {
		router := gin.New()
		router.Use(Metric())

		before := getCounterValue(t, http.MethodGet, "unknown", "404")

		writer := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/missing", nil)
		router.ServeHTTP(writer, request)

		require.Equal(t, http.StatusNotFound, writer.Code)
		after := getCounterValue(t, http.MethodGet, "unknown", "404")
		assert.Equal(t, before+1, after)
	})
}
