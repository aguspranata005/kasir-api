package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Barang struct {
	ID        int    `json:"id"`
	Nama      string `json:"name"`
	Deskripsi string `json:"Deskripsi"`
}

var barang = []Barang{
	{ID: 1, Nama: "Payung", Deskripsi: "Payung merupakan alat untuk melindungi anda dari panas dan hujan."},
	{ID: 2, Nama: "Ember", Deskripsi: "Ember merupakan suatu alat untuk tempat penampung yang diinginkan."},
	{ID: 3, Nama: "Gayung", Deskripsi: "Gayung merupakan sautu alat untuk mengambil sesuatu dari tempat penampungan."},
}

func getBarangByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID produk tidak valid.", http.StatusBadRequest)
		return
	}

	for _, p := range barang {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk tidak ada", http.StatusNotFound)
}

func updateBarang(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Produk ID tidak valid", http.StatusBadRequest)
		return
	}

	var updateBarang Barang
	err = json.NewDecoder(r.Body).Decode(&updateBarang)
	if err != nil {
		http.Error(w, "Permintaan tidak valid", http.StatusBadRequest)
		return
	}

	for i := range barang {
		if barang[i].ID == id {
			updateBarang.ID = id
			barang[i] = updateBarang

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateBarang)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func deleteBarang(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Produk ID tidak valid", http.StatusBadRequest)
		return
	}

	for i, p := range barang {
		if p.ID == id {
			barang = append(barang[:i], barang[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses dihapus",
			})
			return
		}
	}

	http.Error(w, "Produk tidak ada", http.StatusNotFound)
}

func main() {

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getBarangByID(w, r)
		} else if r.Method == "PUT" {
			updateBarang(w, r)
		} else if r.Method == "DELETE" {
			deleteBarang(w, r)
		}

	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(barang)
		} else if r.Method == "POST" {
			var barangBaru Barang
			err := json.NewDecoder(r.Body).Decode(&barangBaru)
			if err != nil {
				http.Error(w, "Request tidak valid", http.StatusBadRequest)
				return
			}

			barangBaru.ID = len(barang) + 1
			barang = append(barang, barangBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses diupdate",
			})
		}
	})

	fmt.Println("Server berjalan di localhost:8081")
	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		fmt.Println("Gagal jalankan server")
	}
}
