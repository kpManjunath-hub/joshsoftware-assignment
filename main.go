package main

import (
	"log"
	"net/http"
)

var (
	productService ProductService
	orderService   OrderService
)

func main() {
	productService = ProductService{
		Products: map[string]Product{
			"product1": {Name: "Product 1", Availability: 100, Price: 10.0, Category: "Premium"},
			"product2": {Name: "Product 2", Availability: 50, Price: 5.0, Category: "Regular"},
			"product3": {Name: "Product 3", Availability: 60, Price: 2.0, Category: "Budget"},
		},
	}
	orderService = OrderService{
		Orders: make(map[string]Order),
	}

	http.HandleFunc("/products", productService.GetProductCatalogue)
	http.HandleFunc("/orders", orderService.PlaceOrder)
	http.HandleFunc("/orders/update", orderService.UpdateOrderStatus)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
