package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phucnguyen/qrify/internal/handlers"
	"github.com/phucnguyen/qrify/internal/models"
	"github.com/phucnguyen/qrify/internal/services"
)

func TestGetQRCodeButNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.GET("/v1/qr/:id", handler.GetQRCode)

	req, _ := http.NewRequest("GET", "/v1/qr/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.URL != "" {
		t.Errorf("Expected URL to be empty, got %s", response.URL)
	}

}

func TestGetQRCodeWithExpiredTime(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.GET("/v1/qr/:id", handler.GetQRCode)

	qr := &models.QRCode{
		ID:        "test",
		URL:       "https://example.com",
		ExpiresAt: time.Now().Add(-time.Hour * 24),
	}
	store.Save(qr)

	req, _ := http.NewRequest("GET", "/v1/qr/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.ExpiresAt.Before(time.Now()) {
		t.Errorf("Expected ExpiresAt to be in the past")
	}
}

func TestGetQRCodeWithValidId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.GET("/v1/qr/:id", handler.GetQRCode)

	qr := &models.QRCode{
		ID:        "test",
		URL:       "https://example.com",
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	store.Save(qr)

	req, _ := http.NewRequest("GET", "/v1/qr/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.URL != qr.URL {
		t.Errorf("Expected URL %s, got %s", qr.URL, response.URL)
	}

}

func TestGenerateQRCodeWithValidURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.POST("/v1/qr", handler.CreateQRCode)

	qr := &models.QRCodeRequest{
		URL:          "https://example.com",
		ExpiresInSec: 3600,
	}

	jsonBody, _ := json.Marshal(qr)

	req, _ := http.NewRequest("POST", "/v1/qr", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.URL != qr.URL {
		t.Errorf("Expected URL %s, got %s", qr.URL, response.URL)
	}

}

func TestGenerateQRCodeWithExpiresInSec(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.POST("/v1/qr", handler.CreateQRCode)

	qr := &models.QRCodeRequest{
		URL:          "https://example.com",
		ExpiresInSec: 3600,
	}

	jsonBody, _ := json.Marshal(qr)

	req, _ := http.NewRequest("POST", "/v1/qr", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	storedQR, err := store.FindByID(response.ID)
	if err != nil {
		t.Fatalf("Failed to get QR code from store: %v", err)
	}

	expectedExpiry := time.Now().Add(time.Hour)
	if diff := storedQR.ExpiresAt.Sub(expectedExpiry); diff < -5*time.Second || diff > 5*time.Second {
		t.Errorf("Expected ExpiresAt to be approximately 1 hour from now, got %v", storedQR.ExpiresAt)
	}
}

func TestGenerateQRCodeWithoutExpiresInSec(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.POST("/v1/qr", handler.CreateQRCode)

	qr := &models.QRCodeRequest{
		URL: "https://example.com",
	}

	jsonBody, _ := json.Marshal(qr)

	req, _ := http.NewRequest("POST", "/v1/qr", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	storedQR, err := store.FindByID(response.ID)
	if err != nil {
		t.Fatalf("Failed to get QR code from store: %v", err)
	}

	if !storedQR.ExpiresAt.IsZero() {
		t.Errorf("Expected ExpiresAt to be zero, got %v", storedQR.ExpiresAt)
	}
}

func TestGenerateQRCodeWithInvalidURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.POST("/v1/qr", handler.CreateQRCode)

	qr := &models.QRCodeRequest{
		URL: "invalid-url",
	}

	jsonBody, _ := json.Marshal(qr)

	req, _ := http.NewRequest("POST", "/v1/qr", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
}

func TestGetQrCodeByUrlWithValidUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.GET("/v1/qr", handler.GetQRCodeByURL)

	qr := &models.QRCode{
		ID:        "test123",
		URL:       "https://example.com",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err := store.Save(qr)
	if err != nil {
		t.Fatalf("Failed to save QR code: %v", err)
	}

	req, _ := http.NewRequest("GET", "/v1/qr?url=https://example.com", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.URL != qr.URL {
		t.Errorf("Expected URL %s, got %s", qr.URL, response.URL)
	}
}

func TestGetQrCodeByUrlWithValidUrlButNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.GET("/v1/qr", handler.GetQRCodeByURL)

	req, _ := http.NewRequest("GET", "/v1/qr?url=https://example.com", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}

	var response models.QRCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.URL != "" {
		t.Errorf("Expected URL to be empty, got %s", response.URL)
	}
}

func TestGetScanCountWithValidId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.GET("/v1/qr/:id/scans", handler.GetScanCount)

	qr := &models.QRCode{
		ID:        "test123",
		URL:       "https://example.com",
		ScanCount: 5,
	}
	err := store.Save(qr)
	if err != nil {
		t.Fatalf("Failed to save QR code: %v", err)
	}

	req, _ := http.NewRequest("GET", "/v1/qr/test123/scans", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response struct {
		ScanCount int `json:"scan_count"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.ScanCount != 5 {
		t.Errorf("Expected scan count 5, got %d", response.ScanCount)
	}
}

func TestGetScanCountWithInvalidId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.GET("/v1/qr/:id/scans", handler.GetScanCount)

	req, _ := http.NewRequest("GET", "/v1/qr/nonexistent/scans", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}

	var response struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Error == "" {
		t.Error("Expected error message, got empty string")
	}
}

func TestGetScanCountAfterIncrement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	store := NewMockQRCodeStore()
	qrService := services.NewQRService(store)
	handler := handlers.NewQRHandler(qrService)
	router.GET("/v1/qr/:id/scans", handler.GetScanCount)

	qr := &models.QRCode{
		ID:        "test123",
		URL:       "https://example.com",
		ScanCount: 0,
	}
	err := store.Save(qr)
	if err != nil {
		t.Fatalf("Failed to save QR code: %v", err)
	}

	err = store.IncrementScanCount("test123")
	if err != nil {
		t.Fatalf("Failed to increment scan count: %v", err)
	}

	req, _ := http.NewRequest("GET", "/v1/qr/test123/scans", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response struct {
		ScanCount int `json:"scan_count"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.ScanCount != 1 {
		t.Errorf("Expected scan count 1, got %d", response.ScanCount)
	}
}
