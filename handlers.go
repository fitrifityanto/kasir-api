package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func handleProduk(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		response := Response{
			Message: "Berhasil mengambil data produk",
			Data:    produk,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case "POST":
		// baca dari Request
		var produkBaru Produk
		err := json.NewDecoder(r.Body).Decode(&produkBaru)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		//masukkan data kedalam variable produk
		produkBaru.ID = len(produk) + 1
		produk = append(produk, produkBaru)

		response := Response{
			Message: "produk berhasil ditambahkan",
			Data:    produkBaru,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			response := Response{
				Message: "Detail produk ditemukan",
				Data:    p,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

func updateProdukByID(w http.ResponseWriter, r *http.Request) {
	//get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// ganti jadi int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	//get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//loop produk, cari id, ganti sesuai data dari request
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			response := Response{
				Message: "Produk berhasil di update",
				Data:    updateProduk,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

func deleteProdukByID(w http.ResponseWriter, r *http.Request) {

	//get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// ganti jadi int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	//loop produk, cari id dan index yang mau dihapus
	for i, p := range produk {
		if p.ID == id {
			// buat slice baru dengan data sebelum dan sesudah index
			produk = append(produk[:i], produk[i+1:]...)

			response := Response{
				Message: "Produk berhasil di hapus",
				Data:    nil,
			}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(response)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}
