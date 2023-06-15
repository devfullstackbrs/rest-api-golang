package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Product representa um produto
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var products []Product
var currentID int

// getAllProducts retorna todos os produtos
func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// getProduct retorna um produto pelo ID
func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	for _, product := range products {
		if product.ID == id {
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	json.NewEncoder(w).Encode(&Product{})
}

// createProduct cria um novo produto
func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	currentID++
	product.ID = currentID
	products = append(products, product)
	json.NewEncoder(w).Encode(product)
}

// updateProduct atualiza um produto existente
func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	for index, product := range products {
		if product.ID == id {
			products = append(products[:index], products[index+1:]...)

			var updatedProduct Product
			if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			updatedProduct.ID = id
			products = append(products, updatedProduct)

			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}

	json.NewEncoder(w).Encode(&Product{})
}

// deleteProduct deleta um produto existente
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	for index, product := range products {
		if product.ID == id {
			products = append(products[:index], products[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(products)
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rotas
	r.Get("/products", getAllProducts)
	r.Get("/products/{id}", getProduct)
	r.Post("/products", createProduct)
	r.Put("/products/{id}", updateProduct)
	r.Delete("/products/{id}", deleteProduct)

	http.ListenAndServe(":8080", r)

}
