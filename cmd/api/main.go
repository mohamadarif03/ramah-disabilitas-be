package main

import (
	"log"
	"os"
	"ramah-disabilitas-be/internal/router"
	"ramah-disabilitas-be/pkg/database"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// 1. Ambil port dari Environment Variable yang diberikan Koyeb
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default jika jalan di lokal
	}

	database.Connect()
	database.Migrate()
	database.SeedAdmin()

	r := router.SetupRouter()

	// 2. WAJIB: Bind ke 0.0.0.0, bukan localhost
	log.Println("Server jalan di port:", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
