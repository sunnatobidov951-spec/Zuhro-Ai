package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"zuhroai/internal/repository/postgres"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL is not set")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	productRepo := postgres.NewProductRepository(db)

	_ = productRepo

	log.Println("Zuhro.AI server started")
}

