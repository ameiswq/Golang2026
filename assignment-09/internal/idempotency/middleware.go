package idempotency

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

func Middleware(store Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}
		cached, exists, err := store.Get(r.Context(), key)
		if err != nil {
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}

		if exists {
			if cached.Completed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(cached.StatusCode)
				w.Write([]byte(cached.Body))
				return
			}

			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}
		started, err := store.StartProcessing(r.Context(), key)
		if err != nil {
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}
		
		if !started {
			cached, exists, err := store.Get(r.Context(), key)
			if err != nil {
				http.Error(w, "storage error", http.StatusInternalServerError)
				return
			}

			if exists && cached.Completed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(cached.StatusCode)
				w.Write([]byte(cached.Body))
				return
			}

			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		body := recorder.Body.String()
		err = store.Finish(r.Context(), key, recorder.Code, body)
		if err != nil {
			http.Error(w, "storage error", http.StatusInternalServerError)
			return
		}

		for name, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		w.WriteHeader(recorder.Code)
		w.Write(bytes.TrimSpace([]byte(body)))
	})
}
