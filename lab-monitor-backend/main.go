package main

import (
	"lab-monitor-backend/db"
	"lab-monitor-backend/handlers"
	"lab-monitor-backend/scheduler"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Baca file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal baca file .env")
	}

	// Konek ke database + bikin tabel
	db.Connect()
	db.Migrate()

	// Jalanin ping scheduler di background
	// "go" artinya jalan bersamaan, ga ngeblock program utama
	go scheduler.StartPing()

	// Setup router API
	r := gin.Default()

	// Routes
	api := r.Group("/api")
	{
		api.GET("/labs", handlers.GetLabs)
		api.POST("/labs", handlers.CreateLab)

		api.GET("/devices", handlers.GetDevices)
		api.POST("/devices", handlers.CreateDevice)
		api.DELETE("/devices/:id", handlers.DeleteDevice)
		api.GET("/scan/stream", handlers.ScanStream)
		api.POST("/scan", handlers.TriggerScan)
	}

	// Jalanin server
	port := os.Getenv("PORT")
	log.Println("Server jalan di port", port)
	r.Run(":" + port)
}