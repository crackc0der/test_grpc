package product

import (
	"context"
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"

	product_grpc "server/api/note_v1"
)

type RepositoryProduct struct {
	db *sqlx.DB
}

func NewRepositoryProduct(db *sqlx.DB) *RepositoryProduct {
	return &RepositoryProduct{db: db}
}

func (r *RepositoryProduct) SelectProducts(_ context.Context) (*product_grpc.AllProductMessage, error) {
	query := "SELECT * FROM product"

	var products *product_grpc.AllProductMessage

	err := r.db.Select(products.GetProducts, query)
	if err != nil {
		return nil, fmt.Errorf("error selecting product in repository's method SelectAllProducts: %w", err)
	}

	return products, nil
}

func (r *RepositoryProduct) SelectProductByID(_ context.Context, id *product_grpc.ProductRequest) (*product_grpc.ProductMessage, error) {
	query := "SELECT * FROM product WHERE product_id=$1"

	var product *product_grpc.ProductMessage

	res := strconv.FormatInt(id.GetId(), 10)

	err := r.db.Get(product, query, res)
	if err != nil {
		return nil, fmt.Errorf("error selecting product in repository's method SelectProductById: %w", err)
	}

	return product, nil
}

func (r *RepositoryProduct) InsertProduct(_ context.Context, prod *product_grpc.ProductMessage) (*product_grpc.ProductMessage, error) {
	query := `INSERT INTO product (product_name, product_category, product_price) VALUES(:ProductName,:ProductCategory,:ProductPrice)`

	_, err := r.db.NamedExec(query, prod)
	if err != nil {
		return nil, fmt.Errorf("error inserting product in repository's method InsertProduct: %w", err)
	}

	return prod, nil
}

func (r *RepositoryProduct) DeleteProductByID(_ context.Context, productID *product_grpc.ProductRequest) (*product_grpc.ProductResponse, error) {
	query := "DELETE * FROM product WHERE product_id=:id"

	prodID := strconv.FormatInt(productID.GetId(), 10)

	_, err := r.db.Exec(query, prodID)
	if err != nil {
		return nil, fmt.Errorf("error deleting product in repository's mothod DeleteProductById: %w", err)
	}

	return &product_grpc.ProductResponse{Deleted: true}, nil
}

func (r *RepositoryProduct) UpdateProduct(_ context.Context, product *product_grpc.ProductMessage) (*product_grpc.ProductMessage, error) {
	query := `UPDATE product SET product_id=:Id product_name=:Product_name product_category=:ProductCategory productPrice=:Price 
			RETURNING id, ProductName, ProductCategory, Price`

	var updatedProduct *product_grpc.ProductMessage

	err := r.db.QueryRowx(query, product).StructScan(updatedProduct)
	if err != nil {
		return nil, fmt.Errorf("error updating product in repository's method UpdateProduct: %w", err)
	}

	return updatedProduct, nil
}
