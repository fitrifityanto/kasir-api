package main

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Beras", Harga: 70000, Stok: 10},
	{ID: 2, Nama: "Gula", Harga: 15000, Stok: 20},
	{ID: 3, Nama: "Kopi Kacamata", Harga: 12000, Stok: 5},
}
