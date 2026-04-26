package demo

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"assignment-09/internal/idempotency"
	"assignment-09/internal/payment"
	"sync"

	_ "github.com/lib/pq"
)

func RunIdempotencyDemo() {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=practice9 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("database open error:", err)
		return
	}
	defer db.Close()
	clearOldDemoData(db)
	store := idempotency.NewPostgresStore(db)
	mainHandler := http.HandlerFunc(payment.RepaymentHandler)
	server := httptest.NewServer(idempotency.Middleware(store, mainHandler))
	defer server.Close()
	key := "demo-payment-key-1"
	var wg sync.WaitGroup
	fmt.Println("Sending 10 simultaneous requests with the same Idempotency-Key")
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(requestNumber int) {
			defer wg.Done()
			sendPaymentRequest(server.URL, key, requestNumber)
		}(i)
	}
	wg.Wait()

	fmt.Println("Sending one more request after completion with the same key")
	sendPaymentRequest(server.URL, key, 11)
}


func clearOldDemoData(db *sql.DB) {
	_, err := db.Exec(`DELETE FROM idempotency_keys WHERE key = 'demo-payment-key-1'`)
	if err != nil {
		fmt.Println("cleanup error:", err)
	}
}

func sendPaymentRequest(url string, key string, requestNumber int) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Printf("Request %d error: %v\n", requestNumber, err)
		return
	}

	request.Header.Set("Idempotency-Key", key)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Request %d error: %v\n", requestNumber, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Request %d: status=%d body=%s\n", requestNumber, resp.StatusCode, string(body))
}
