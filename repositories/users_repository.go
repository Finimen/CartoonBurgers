package repositories

import (
	"CartoonBurgers/models"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type UserRepository struct {
	db *sql.DB
}

func (repo *UserRepository) Init(ctx context.Context, db *sql.DB) error {
	var path = filepath.Join("..", "repositories", "migrations", "001_create_users_table_up.sql")
	var com, err = os.ReadFile(path)
	if err != nil {
		return err
	}

	repo.db = db
	_, err = db.ExecContext(ctx, string(com))
	if err != nil {
		fmt.Print("EXEC ERROR", path)
		return err
	}

	return nil
}

func (r *UserRepository) GetUserProfile(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	username = strings.Replace(username, " ", "", -1)
	query := `SELECT username, email, bonus FROM users WHERE username = $1`
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.Username,
		&user.Email,
		&user.Bonus,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetUserByUsername(ctx context.Context, name string) (string, error) {
	var password string
	var row = repo.db.QueryRowContext(ctx, "SELECT passwordHash FROM users WHERE username = ?", name)
	err := row.Scan(&password)
	return password, err
}

func (repo *UserRepository) CreateUser(ctx context.Context, user models.User, hashedPassword string) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (username, passwordHash, email) VALUES (?, ?, ?)", user.Username, hashedPassword, user.Email)
	return err
}
