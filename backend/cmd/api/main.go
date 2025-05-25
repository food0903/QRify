package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/phucnguyen/qrify/internal/handlers"
	"github.com/phucnguyen/qrify/internal/services"
)

func main() {
	db, err := sql.Open("postgres", "postgres://qrify_user:postgres@localhost:5432/qrify?sslmode=disable")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    store := services.NewPostgresQRCodeStore(db)
    qrService := services.NewQRService(store)

	qrHandler := handlers.NewQRHandler(qrService)

	r := gin.Default()


	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// QR code endpoints
	qr := r.Group("/v1/qr")
	{
		qr.POST("", qrHandler.CreateQRCode)
		qr.GET("/:id", qrHandler.GetQRCode)
		qr.DELETE("/:id", qrHandler.DeleteQRCode)
	}

	// Redirect endpoint for QR code scans
	r.GET("/r/:id", qrHandler.HandleRedirect)

	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 