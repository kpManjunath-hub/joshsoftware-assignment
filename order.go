package main

import (
	"encoding/json"
	"net/http"
)

type OrderService struct {
	Orders map[string]Order
}

type Order struct {
	OrderID      string  `json:"orderID"`
	ProductID    string  `json:"productID"`
	Quantity     int     `json:"quantity"`
	OrderValue   float64 `json:"orderValue"`
	DispatchDate string  `json:"dispatchDate,omitempty"`
	OrderStatus  string  `json:"orderStatus"`
	Discounted   bool    `json:"-"`
}

func (o *OrderService) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, ok := productService.Products[order.ProductID]
	if !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if order.Quantity < 1 || order.Quantity > 10 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	orderValue := product.Price * float64(order.Quantity)

	hasPremium := false
	for _, existingOrder := range o.Orders {
		existingProduct, exists := productService.Products[existingOrder.ProductID]
		if exists && existingProduct.Category == "Premium" {
			hasPremium = true
			break
		}
	}

	if product.Category == "Premium" && !hasPremium {
		orderValue *= 0.9 // Apply 10% discount
		order.Discounted = true
	}

	order.OrderValue = orderValue
	order.OrderStatus = "Placed"
	o.Orders[order.OrderID] = order

	product.Availability -= order.Quantity
	productService.Products[order.ProductID] = product

	jsonData, _ := json.Marshal(order)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (o *OrderService) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderID")
	if orderID == "" {
		http.Error(w, "Missing orderID parameter", http.StatusBadRequest)
		return
	}

	order, ok := o.Orders[orderID]
	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	var payload struct {
		OrderStatus  string `json:"orderStatus"`
		DispatchDate string `json:"dispatchDate,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.OrderStatus = "Dispatched"
	order.DispatchDate = payload.DispatchDate
	o.Orders[orderID] = order

	jsonData, _ := json.Marshal(order)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
