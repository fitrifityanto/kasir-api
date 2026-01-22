package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getProdukByID(w, r)
		case "PUT":
			updateProdukByID(w, r)
		case "DELETE":
			deleteProdukByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		}
	})

	// route produk tanpa ID
	http.HandleFunc("/api/produk", handleProduk)

	// route category
	http.HandleFunc("/api/categories", handleCategories)

	// route category dengan ID
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			getCategoryByID(w, r)
		case "PUT":
			updateCategoryByID(w, r)
		case "DELETE":
			deleteCategoryByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "API running"})
	})
	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running sever")
	}
}
