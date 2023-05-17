package main

import (
	"encoding/json"
	"net/http"
)

type ProductService struct {
	Products map[string]Product
}

type Product struct {
	Name         string  `json:"name"`
	Availability int     `json:"availability"`
	Price        float64 `json:"price"`
	Category     string  `json:"category"`
}

func (p *ProductService) GetProductCatalogue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(p.Products); err != nil {
		http.Error(w, "Failed to encode product catalogue", http.StatusInternalServerError)
		return
	}
}
