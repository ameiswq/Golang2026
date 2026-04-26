package exchangeService

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := RateResponse{
				Base:   "USD",
				Target: "EUR",
				Rate:   0.9,
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		service := NewExchangeService(server.URL)
		rate, err := service.GetRate("USD", "EUR")
		assert.NoError(t, err)
		assert.Equal(t, 0.9, rate)
	})

	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid currency pair",
			})
		}))
		defer server.Close()
		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "XXX")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "api error")
	})

	t.Run("malformed json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()
		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decode error")
	})

	t.Run("empty body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")
		assert.Error(t, err)
	})

	t.Run("server error 500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()
		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status")
	})
}