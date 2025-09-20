package repositories

import (
	"CartoonBurgers/models"
	"context"
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type MenuPerository struct {
	db *sql.DB
}

func (prod *MenuPerository) Init(ctx context.Context) error {
	var db, err = sql.Open("sqlite", "./products")
	if err != nil {
		return err
	}

	var path = filepath.Join("migrations", "001_create_products_table_up.sql")
	var req []byte
	req, err = os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(req))

	prod.db = db

	return err
}

func (prod *MenuPerository) GetAll(ctx context.Context) ([]models.Product, error) {
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
