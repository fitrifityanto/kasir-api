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
