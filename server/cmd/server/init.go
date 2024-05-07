package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	product_grpc "server/api/note_v1"

	"server/config"
	"server/internal/category"
	"server/internal/product"
)

func Run() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	dsn := config.GetDsn()

	dbConn, err := sqlx.Connect(config.DataBase.DBType, dsn)
	if err != nil {
		log.Fatal(err)
	}

	productRepository := product.NewRepositoryProduct(dbConn)
	productService := product.NewServiceProduct(productRepository)
	productEndpoint := product.NewEndpointProduct(productService, logger)

	categoryRepository := category.NewRepositoryCategory(dbConn)
	categoryService := category.NewServiceCategory(categoryRepository)
	categoryEndpoint := category.NewEndpointCategory(categoryService, logger)

	conn, err := net.Listen("tcp", config.Host.HostPort)
	if err != nil {
		log.Fatal("error in start grpc server: %w", err)
	}

	serv := grpc.NewServer()

	reflection.Register(serv)
	product_grpc.RegisterProductServer(serv, productEndpoint)
	product_grpc.RegisterCategoryServer(serv, categoryEndpoint)

	if err := serv.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
