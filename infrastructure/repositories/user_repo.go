package repositories

import (
	"GoSyntaxDoc/domain/entities"
	"GoSyntaxDoc/infrastructure/database"
	"GoSyntaxDoc/presentation/middleware"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// ✅ UserRepository Struct (Using pgxpool.Pool)
type UserRepository struct {
	DB *pgxpool.Pool
}

// ✅ NewUserRepository - Accepts `DBinstance` for proper dependency injection
func NewUserRepository(db *database.DBinstance) *UserRepository {
	return &UserRepository{DB: db.DB}
}

// ✅ Fetch User by ID (Updated for pgx)
func (repo *UserRepository) FetchUserById(userId int) (*entities.User, error) {
	ctx := context.Background()
	query := `SELECT id, first_name, last_name, created_at FROM users WHERE id = $1`

	var user entities.User
	var createdAt pgtype.Timestamp

	err := repo.DB.QueryRow(ctx, query, userId).Scan(
		&user.ID, &user.FirstName, &user.LastName, &createdAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			middleware.Log.WithFields(logrus.Fields{"user_id": userId}).Warn("User not found")
			return nil, err
		}
		middleware.Log.WithFields(logrus.Fields{"error": err}).Error("Database Query Error")
		return nil, err
	}

	// ✅ Convert NULL timestamps to zero-value JSONTime
	if createdAt.Valid {
		user.CreatedAt = entities.JSONTime{Time: createdAt.Time}
	} else {
		user.CreatedAt = entities.JSONTime{} // Empty timestamp
	}

	return &user, nil
}

// ✅ Create User (Updated for pgx)
func (repo *UserRepository) CreateUser(first_name string, last_name string) (*entities.User, error) {
	ctx := context.Background()
	query := `INSERT INTO users (first_name, last_name) VALUES ($1, $2) RETURNING id, created_at`

	var user entities.User
	var createdAt pgtype.Timestamp // ✅ Use pgx.NullTime instead of sql.NullTime

	err := repo.DB.QueryRow(ctx, query, first_name, last_name).Scan(
		&user.ID, &createdAt,
	)

	if err != nil {
		middleware.Log.WithFields(logrus.Fields{"error": err}).Error("Database Query Error")
		return nil, err
	}

	// ✅ Convert NULL timestamps to zero-value JSONTime
	if createdAt.Valid {
		user.CreatedAt = entities.JSONTime{Time: createdAt.Time}
	} else {
		user.CreatedAt = entities.JSONTime{} // Empty timestamp
	}

	user.FirstName = first_name
	user.LastName = last_name

	return &user, nil
}
