package models

import "time"

type QRCode struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	ImageBase64 string    `json:"image_base64,omitempty"`
	ScanCount   int       `json:"scan_count"`
}

type QRCodeRequest struct {
	URL          string `json:"url" binding:"required,url"`
	ExpiresInSec int64  `json:"expires_in_sec,omitempty"`
}

type QRCodeResponse struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	QRCodeURL   string    `json:"qr_code_url"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	ImageBase64 string    `json:"image_base64,omitempty"`
	ScanCount   int       `json:"scan_count"`
}
