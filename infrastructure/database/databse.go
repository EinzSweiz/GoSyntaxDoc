package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBinstance struct {
	DB *pgxpool.Pool
}

var Database DBinstance

func ConnectDB() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		os.Exit(2)
	}
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatal("Error creating connection pool: ", err)
		os.Exit(2)
	}
	Database = DBinstance{DB: db}

}

func CloseDB() {
	if Database.DB != nil {
		Database.DB.Close()
		log.Println("Database connection closed")
	}
}
