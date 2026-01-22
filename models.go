package main

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var produk = []Produk{
	{ID: 1, Nama: "Beras", Harga: 70000, Stok: 10},
	{ID: 2, Nama: "Gula", Harga: 15000, Stok: 20},
	{ID: 3, Nama: "Kopi Kacamata", Harga: 12000, Stok: 5},
}

var categories = []Category{
	{ID: 1, Name: "Sembako", Description: "Kategori produk sembako"},
	{ID: 2, Name: "Perawatan Tubuh", Description: "Kategori produk perawatan tubuh"},
	{ID: 3, Name: "Pembersih Rumah", Description: "Kategori produk pembersih rumah"},
}
