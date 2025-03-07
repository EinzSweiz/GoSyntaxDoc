package migrations

import (
	"GoSyntaxDoc/infrastructure/database"
	"context"
	"log"
)

func InitOrder() {
	ctx := context.Background()

	_, err := database.Database.DB.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id),
            product_id INTEGER REFERENCES products(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)

	if err != nil {
		log.Fatalf("Error creating orders table: %v", err)
	}
	log.Println("Orders table created successfully")
}
