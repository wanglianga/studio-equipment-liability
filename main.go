package main

import (
	"fmt"
	"log"
	"net/http"

	"studio-equipment-manager/handler"
	"studio-equipment-manager/service"
	"studio-equipment-manager/store"
)

func main() {
	s := store.New()
	svc := service.New(s)
	h := handler.New(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	mux.HandleFunc("GET /health", h.HealthCheck)

	addr := ":8090"
	fmt.Printf("Studio Equipment Manager starting on %s\n", addr)
	fmt.Println("API endpoints:")
	fmt.Println("  POST   /api/equipments            - Create equipment")
	fmt.Println("  GET    /api/equipments            - List equipment")
	fmt.Println("  GET    /api/equipments/{id}       - Get equipment")
	fmt.Println("  POST   /api/borrow                - Borrow equipment")
	fmt.Println("  POST   /api/borrow/return         - Return & inspect")
	fmt.Println("  GET    /api/borrow                - List borrow records")
	fmt.Println("  GET    /api/borrow/{id}           - Get borrow record")
	fmt.Println("  POST   /api/damage                - Register damage")
	fmt.Println("  GET    /api/damage                - List damage reports")
	fmt.Println("  GET    /api/damage/{id}           - Get damage report")
	fmt.Println("  POST   /api/repair-quote          - Create repair quote")
	fmt.Println("  POST   /api/repair-complete/{id}  - Complete repair")
	fmt.Println("  POST   /api/deduction             - Deduct deposit")
	fmt.Println("  POST   /api/deduction/accessory   - Deduct accessory")
	fmt.Println("  POST   /api/accessory-prices      - Add accessory price")
	fmt.Println("  GET    /api/accessory-prices      - List accessory prices")
	fmt.Println("  POST   /api/appeal                - Create appeal")
	fmt.Println("  POST   /api/appeal/review         - Review appeal")
	fmt.Println("  GET    /api/appeal                - List appeals")
	fmt.Println("  GET    /health                    - Health check")

	log.Fatal(http.ListenAndServe(addr, mux))
}
