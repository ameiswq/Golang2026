package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PaymentResponse struct {
	Status string `json:"status"`
	Amount int `json:"amount"`
	TransactionID string `json:"transaction_id"`
}

func RepaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Processing started")
	time.Sleep(2 * time.Second)

	response := PaymentResponse{
		Status: "paid",
		Amount: 1000,
		TransactionID: "uuid-demo-12345",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	fmt.Println("Processing completed")
}
