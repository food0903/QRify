package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phucnguyen/qrify/internal/models"
	"github.com/phucnguyen/qrify/internal/services"
)

type QRHandler struct {
	qrService *services.QRService
}

func NewQRHandler(qrService *services.QRService) *QRHandler {
	return &QRHandler{
		qrService: qrService,
	}
}

// create the qr code by url
func (h *QRHandler) CreateQRCode(c *gin.Context) {
	var req models.QRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qr, err := h.qrService.GenerateQRCode(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, qr)
}

// get the qr code in base64 format
func (h *QRHandler) GetQRCode(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("DEBUG: id param =", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code ID is required"})
		return
	}

	qr, err := h.qrService.GetQRCode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if qr == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR code not found"})
		return
	}

	c.JSON(http.StatusOK, qr)
}

// delete the qr code by id
func (h *QRHandler) DeleteQRCode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code ID is required"})
		return
	}

	err := h.qrService.DeleteQRCode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
