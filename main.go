package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middlewares"

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
	APIKey string `mapstructure:"API_KEY"`
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
		APIKey: viper.GetString("API_KEY"),
	}

	// setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()

	apiKeyMiddleware := middlewares.APIKey(config.APIKey)

	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	reportRepo := repositories.NewReportRepository(db)

	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productService := services.NewProductService(productRepo, categoryRepo)
	productHandler := handlers.NewProductHandler(productService)

	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/product", apiKeyMiddleware(productHandler.HandleProducts))
	mux.HandleFunc("/api/product/", apiKeyMiddleware(productHandler.HandleProductByID))

	mux.HandleFunc("/api/categories", categoryHandler.HandleCategory)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	mux.HandleFunc("/api/checkout", apiKeyMiddleware(transactionHandler.HandleCheckout))

	mux.HandleFunc("/api/report/hari-ini", apiKeyMiddleware(reportHandler.HandleReport))
	mux.HandleFunc("/api/report", apiKeyMiddleware(reportHandler.GetReport))

	// Route umum
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "API running"})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"app": "Kasir API app", "maintainer": "Fitriningtyas", "message": "Server is running"})
	})

	handlerWithCORS := middlewares.CORS(mux)

	fmt.Println("Server running di localhost:" + config.Port)

	err = http.ListenAndServe(":"+config.Port, handlerWithCORS)
	if err != nil {
		fmt.Println("gagal running server:", err)
	}
}
