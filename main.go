package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// --- 1. STRUKTUR DATA ---
type Product struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
	Img      string `json:"img"`
}

type LoginRequest struct {
	Password string `json:"password"`
}

// --- 2. DATABASE SEMENTARA (IN-MEMORY) ---
var (
	products []Product
	mutex    sync.Mutex // Agar aman saat diakses banyak orang
	idCounter int64 = 100
)

func init() {
	// Data Awal (Seed Data)
	products = []Product{
		{ID: 1, Name: "Nugget Kanzler", Category: "Nugget", Price: 45000, Img: "https://images.unsplash.com/photo-1569691105751-88df003de7a4?w=500&q=80"},
		{ID: 2, Name: "Sosis Kanzler Beef", Category: "Sosis", Price: 48000, Img: "https://images.unsplash.com/photo-1585325701165-351af92f9656?w=500&q=80"},
		{ID: 3, Name: "Daging Slice 500g", Category: "Daging", Price: 65000, Img: "https://images.unsplash.com/photo-1607623814075-e51df1bdc82f?w=500&q=80"},
        {ID: 4, Name: "Kentang Goreng", Category: "Snack", Price: 28000, Img: "https://images.unsplash.com/photo-1630384060421-cb20d0e0649d?w=500&q=80"},
	}
}

// --- 3. HANDLERS (FUNGSI SERVER) ---

// Menampilkan Halaman Utama
func viewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	tmpl.Execute(w, nil)
}

// API: Ambil Semua Produk (GET)
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	json.NewEncoder(w).Encode(products)
	mutex.Unlock()
}

// API: Tambah/Edit Produk (POST)
func saveProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if p.ID == 0 {
		// Create Baru
		idCounter++
		p.ID = idCounter
		products = append(products, p)
	} else {
		// Update Lama
		for i, item := range products {
			if item.ID == p.ID {
				products[i] = p
				break
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

// API: Hapus Produk (DELETE)
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	mutex.Lock()
	defer mutex.Unlock()

	newProducts := []Product{}
	for _, p := range products {
		if p.ID != id {
			newProducts = append(newProducts, p)
		}
	}
	products = newProducts
	w.WriteHeader(http.StatusOK)
}

// API: Login Admin Sederhana
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	// PASSWORD ADMIN DISINI
	if req.Password == "admin123" {
		// Set Cookie sederhana untuk tanda login
		http.SetCookie(w, &http.Cookie{
			Name:    "admin_session",
			Value:   "true",
			Expires: time.Now().Add(24 * time.Hour),
		})
		w.Write([]byte(`{"status":"success"}`))
	} else {
		http.Error(w, "Password Salah", http.StatusUnauthorized)
	}
}

func main() {
	// Routing
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/api/products", getProducts)      // GET list
	http.HandleFunc("/api/product/save", saveProduct)  // POST save
	http.HandleFunc("/api/product/delete", deleteProduct) // DELETE
	http.HandleFunc("/api/login", loginHandler)        // POST login

	// Serve Static Files (Untuk gambar jika ada folder static)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server Aulia Frozen berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
