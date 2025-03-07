package migrations

import (
	"GoSyntaxDoc/infrastructure/database"
	"context"
	"log"
)

func InitProduct() {
	ctx := context.Background()

	_, err := database.Database.DB.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            serial_number VARCHAR(255) UNIQUE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)

	if err != nil {
		log.Fatalf("Error creating products table: %v", err)
	}
	log.Println("Products table created successfully")
}
