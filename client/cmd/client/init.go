package main

import (
	product_grpc "client/api/note_v1"
	"client/config"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"client/internal/category"
	"client/internal/product"
)

func Run() {
	mux := http.NewServeMux()
	timeout := 10

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	conn, err := grpc.Dial(config.GrpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	categoryClient := product_grpc.NewCategoryClient(conn)

	categoryRepository := category.NewRepository(categoryClient)
	categoryService := category.NewService(categoryRepository)
	categoryEndpoint := category.NewEndpoint(categoryService, logger)

	productClient := product_grpc.NewProductClient(conn)

	productRepository := product.NewRepository(productClient)
	productService := product.NewService(productRepository)
	productEndpoint := product.NewEndpoint(productService, logger)

	mux.HandleFunc("GET /get/categories", categoryEndpoint.GetCategories)
	mux.HandleFunc("GET /add/category", categoryEndpoint.AddCategory)
	mux.HandleFunc("GET /update/category", categoryEndpoint.UpdateCategory)
	mux.HandleFunc("GET /delete/category", categoryEndpoint.DeleteCategory)

	mux.HandleFunc("GET /get/products", productEndpoint.GetProducts)
	mux.HandleFunc("GET /get/product", productEndpoint.GetProduct)
	mux.HandleFunc("GET /add/product", productEndpoint.AddProduct)
	mux.HandleFunc("GET /update/product", productEndpoint.UpdateProduct)
	mux.HandleFunc("GET /delete/product", productEndpoint.DeleteProduct)

	srv := http.Server{
		Addr:              config.HttpPort,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * time.Duration(timeout),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
