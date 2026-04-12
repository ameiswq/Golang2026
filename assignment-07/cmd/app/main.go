package main

import (
	"log"
	"fmt"
	"os"
	"assignment-07/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	fmt.Println("JWT_SECRET from env =", os.Getenv("JWT_SECRET"))
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}