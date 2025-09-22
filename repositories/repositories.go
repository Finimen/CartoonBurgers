package repositories

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type AppRepository struct {
	DB *sql.DB

	*UserRepository
	*ProductRerository
}

func NewAppRepository(ctx context.Context, dbPath string) (*AppRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	repo := &AppRepository{DB: db}

	repo.UserRepository = &UserRepository{db: db}
	repo.ProductRerository = &ProductRerository{db: db}

	if err := repo.initUsersTable(ctx); err != nil {
		return nil, err
	}
	if err := repo.initProductsTable(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *AppRepository) initUsersTable(ctx context.Context) error {
	return r.UserRepository.Init(ctx, r.DB)
}

func (r *AppRepository) initProductsTable(ctx context.Context) error {
	return r.ProductRerository.Init(ctx, r.DB)
}
