package migrations

import (
	"GoSyntaxDoc/infrastructure/database"
	"context"
	"log"
)

func InitUser() {
	ctx := context.Background()

	_, err := database.Database.DB.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            first_name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	log.Println("Users table created successfully")
}
