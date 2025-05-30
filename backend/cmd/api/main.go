package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phucnguyen/qrify/internal/database"
	"github.com/phucnguyen/qrify/internal/handlers"
	"github.com/phucnguyen/qrify/internal/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	_ = godotenv.Load(".env.local")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	store := services.NewPostgresQRCodeStore(db)
	qrService := services.NewQRService(store)
	qrHandler := handlers.NewQRHandler(qrService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// QR code endpoints
	qr := r.Group("/v1/qr")
	{
		qr.POST("", qrHandler.CreateQRCode)
		qr.GET("/:id", qrHandler.GetQRCode)
		qr.DELETE("/:id", qrHandler.DeleteQRCode)
		qr.GET("", qrHandler.GetQRCodeByURL)
		qr.GET("/:id/scans", qrHandler.GetScanCount)
	}

	// redirect endpoint for QR code scans
	r.GET("/r/:id", qrHandler.HandleRedirect)

	port := os.Getenv("PORT")

	log.Printf("Server starting on port %s", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
