package main

import (
	"math/rand"
	"os"
	"assignment-09/internal/demo"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	switch os.Args[1] {
	case "retry":
		demo.RunRetryDemo()
	case "idempotency":
		demo.RunIdempotencyDemo()
	}
}
