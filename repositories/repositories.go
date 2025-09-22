package repositories

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type AppRepository struct {
	DB *sql.DB

	*UserRepository
	*MenuRerository
}

func NewAppRepository(ctx context.Context, dbPath string) (*AppRepository, error) {
	fmt.Println("PATH |", dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Print("ERROR CODE 1")
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		fmt.Print("ERROR CODE 2")
		return nil, err
	}

	repo := &AppRepository{DB: db}

	repo.UserRepository = &UserRepository{db: db}
	repo.MenuRerository = &MenuRerository{db: db}

	if err := repo.initUsersTable(ctx); err != nil {
		fmt.Print("ERROR HERE")
		return nil, err
	}
	if err := repo.initProductsTable(ctx); err != nil {
		fmt.Print("ERROR THERE")
		return nil, err
	}

	fmt.Print("ALL WORKS")
	return repo, nil
}

func (r *AppRepository) initUsersTable(ctx context.Context) error {
	return r.UserRepository.Init(ctx, r.DB)
}

func (r *AppRepository) initProductsTable(ctx context.Context) error {
	return r.MenuRerository.Init(ctx, r.DB)
}
