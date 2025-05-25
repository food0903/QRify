package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var qrScansTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "qr_scans_total",
		Help: "Total number of QR code scans",
	},
	[]string{"qr_id"},
)

func init() {
	prometheus.MustRegister(qrScansTotal)
}

// redirect to the server to aggregate the metrics
func (h *QRHandler) HandleRedirect(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code ID is required"})
		return
	}

	qrScansTotal.WithLabelValues(id).Inc()

	qr, err := h.qrService.GetQRCode(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if qr == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR code not found"})
		return
	}

	c.Redirect(http.StatusFound, qr.URL)
} 