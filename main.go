package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// setup routes
	http.HandleFunc("/api/product", productHandler.HandleProducts)

	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	// http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case "GET":
	// 		getProdukByID(w, r)
	// 	case "PUT":
	// 		updateProdukByID(w, r)
	// 	case "DELETE":
	// 		deleteProdukByID(w, r)
	// 	default:
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//
	// 	}
	// })

	// route produk tanpa ID
	// http.HandleFunc("/api/produk", handleProduk)

	// route category
	// http.HandleFunc("/api/categories", handleCategories)

	// route category dengan ID
	// http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
	//
	// 	switch r.Method {
	// 	case "GET":
	// 		getCategoryByID(w, r)
	// 	case "PUT":
	// 		updateCategoryByID(w, r)
	// 	case "DELETE":
	// 		deleteCategoryByID(w, r)
	// 	default:
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	}
	// })

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "API running"})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"app": "Kasir API app", "maintainer": "Fitriningtyas", "message": "Server is running"})
	})

	fmt.Println("Server running di localhost:" + config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("gagal running sever")
	}
}
