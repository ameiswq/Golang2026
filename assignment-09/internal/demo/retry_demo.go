package demo

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"assignment-09/internal/retry"
	"time"
)

func RunRetryDemo() {
	counter := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
		if counter <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error":"payment gateway unavailable"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request, err := http.NewRequest(http.MethodPost, server.URL, nil)
	if err != nil {
		fmt.Println("request error:", err)
		return
	}

	client := retry.PaymentClient{
		Client: &http.Client{},
		MaxRetries: 5,
		BaseDelay: 500 * time.Millisecond,
		MaxDelay: 5 * time.Second,
	}

	body, err := client.ExecutePayment(ctx, request)
	if err != nil {
		fmt.Println("Final error:", err)
		return
	}

	fmt.Println("Final response:", string(body))
}
