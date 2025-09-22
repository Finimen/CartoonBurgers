package repositories

import (
	"CartoonBurgers/models"
	"context"
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type ProductRerository struct {
	db *sql.DB
}

func (prod *ProductRerository) Init(ctx context.Context, db *sql.DB) error {
	var path = filepath.Join("..", "repositories", "migrations", "001_create_products_table_up.sql")
	var req, err = os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(req))

	prod.db = db

	if err != nil {
		return err
	}

	return prod.fillDB(ctx)
}

func (prod *ProductRerository) fillDB(ctx context.Context) error {
	var count int
	var err = prod.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		tx, err := prod.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		products := []models.Product{
			{Name: "Cheese Burger", Price: 160, Count: 1, Type: 0, Category: 0},
			{Name: "Classic Carton", Price: 200, Count: 1, Type: 0, Category: 0},
			{Name: "Twicer Classic", Price: 270, Count: 1, Type: 0, Category: 0},
			{Name: "Twicer Double", Price: 330, Count: 1, Type: 0, Category: 0},
			{Name: "Purple Burger", Price: 390, Count: 1, Type: 1, Category: 0},
			{Name: "Chiken Nuggets", Price: 129, Count: 12, Type: 0, Category: 1},
			{Name: "Efilio Cake", Price: 199, Count: 1, Type: 0, Category: 3},
			{Name: "Purple Cake", Price: 199, Count: 1, Type: 1, Category: 3},
		}
		for _, p := range products {
			_, err = tx.ExecContext(ctx, `INSERT INTO products (pName, pPrice, pCount, pType, pCategory) VALUES (?, ?, ?, ?, ?)`,
				p.Name, p.Price, p.Count, p.Type, p.Category)

			if err != nil {
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return err
}

func (prod *ProductRerository) GetAll(ctx context.Context) ([]models.Product, error) {
	rows, err := prod.db.QueryContext(ctx, "SELECT id, pName, pPrice, pCount, pType, pCategory FROM products")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Count, &p.Type, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, rows.Err()
}
