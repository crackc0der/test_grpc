package category

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	product_grpc "server/api/note_v1"
)

type RepositoryCategory struct {
	db *sqlx.DB
}

func NewRepositoryCategory(db *sqlx.DB) *RepositoryCategory {
	return &RepositoryCategory{db: db}
}

func (r *RepositoryCategory) SelectCategories(_ context.Context) (*product_grpc.AllCategoryMessage, error) {
	query := `SELECT * FROM category`

	var categories []*product_grpc.CategoryMessage

	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("error in repository's method SelectCategories: %w", err)
	}

	allCategories := &product_grpc.AllCategoryMessage{
		Categories: categories,
	}
	return allCategories, nil
}

func (r *RepositoryCategory) InsertCategory(_ context.Context, cat *product_grpc.CategoryMessage) (*product_grpc.CategoryMessage, error) {
	query := `INSERT INTO category (category_name) VALUES ($1) RETURNING category_id, category_name`

	category := &product_grpc.CategoryMessage{}

	err := r.db.QueryRowx(query, cat.CategoryName).StructScan(category)
	if err != nil {
		return nil, fmt.Errorf("error in repository's method InsertCategory: %w", err)
	}

	return category, nil
}

func (r *RepositoryCategory) UpdateCategory(_ context.Context, cat *product_grpc.CategoryMessage) (*product_grpc.CategoryMessage, error) {
	query := `UPDATE category SET category_name=$1 WHERE category_id=$2 RETURNING category_id, category_name`

	category := &product_grpc.CategoryMessage{}

	err := r.db.QueryRowx(query, cat.CategoryName, cat.Id).StructScan(category)
	if err != nil {
		return nil, fmt.Errorf("error in repository's method UpdateCategory: %w", err)
	}

	return category, nil
}

func (r *RepositoryCategory) DeleteCategory(_ context.Context, id *product_grpc.CategoryRequest) (*product_grpc.CategoryResponse, error) {
	query := `DELETE FROM category WHERE category_id=$1`

	_, err := r.db.Exec(query, id.GetId())
	if err != nil {
		return nil, fmt.Errorf("error in repository's method DeleteCategory: %w", err)
	}

	return &product_grpc.CategoryResponse{Deleted: true}, nil
}
