package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"bytes"
	"image/png"

	_ "github.com/lib/pq"
	"github.com/phucnguyen/qrify/internal/models"
	"github.com/skip2/go-qrcode"
)

type QRService struct {
	store QRCodeStore
}

func NewQRService(store QRCodeStore) *QRService {
	return &QRService{
		store: store,
	}
}

// generate qr code by url
func (s *QRService) GenerateQRCode(req *models.QRCodeRequest) (*models.QRCodeResponse, error) {
	if req.URL == "" {
		return nil, errors.New("URL is required")
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	redirectURL := "http://localhost:8080/r/" + id
	qrImg, err := qrcode.New(redirectURL, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrImg.DisableBorder = true
	img := qrImg.Image(256)

	// Add the ID as text below the QR code
	imgWithText, err := addTextBelow(img, id)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, imgWithText); err != nil {
		return nil, err
	}

	expiresAt := time.Time{}
	if req.ExpiresInSec > 0 {
		expiresAt = time.Now().UTC().Add(time.Duration(req.ExpiresInSec) * time.Second)
	}

	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())

	qr := &models.QRCode{
		ID:          id,
		URL:         req.URL,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		ImageBase64: base64Img,
	}

	if err := s.store.Save(qr); err != nil {
		return nil, err
	}

	response := &models.QRCodeResponse{
		ID:          qr.ID,
		URL:         qr.URL,
		QRCodeURL:   "/r/" + qr.ID,
		CreatedAt:   qr.CreatedAt,
		ExpiresAt:   qr.ExpiresAt,
		ImageBase64: qr.ImageBase64,
	}

	return response, nil
}

// get qr code by id
func (s *QRService) GetQRCode(id string) (*models.QRCodeResponse, error) {
	qr, err := s.store.FindByID(id)
	if err != nil {
		return nil, err
	}
	if qr == nil {
		return nil, errors.New("QR code not found")
	}

	return &models.QRCodeResponse{
		ID:          qr.ID,
		URL:         qr.URL,
		QRCodeURL:   "/r/" + qr.ID,
		CreatedAt:   qr.CreatedAt,
		ExpiresAt:   qr.ExpiresAt,
		ImageBase64: qr.ImageBase64,
		ScanCount:   qr.ScanCount,
	}, nil
}

// delete qr code by id
func (s *QRService) DeleteQRCode(id string) error {
	if err := s.store.DeleteByID(id); err != nil {
		return err
	}
	return nil
}

// generate a random ID
func generateID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// get qr code by url
func (s *QRService) GetQRCodeByURL(url string) (*models.QRCodeResponse, error) {
	qr, err := s.store.FindByURL(url)
	if err != nil {
		return nil, err
	}
	if qr == nil {
		return nil, nil
	}
	return &models.QRCodeResponse{
		ID:          qr.ID,
		URL:         qr.URL,
		QRCodeURL:   "/r/" + qr.ID,
		CreatedAt:   qr.CreatedAt,
		ExpiresAt:   qr.ExpiresAt,
		ImageBase64: qr.ImageBase64,
	}, nil
}

func (s *QRService) IncrementScanCount(id string) error {
	return s.store.IncrementScanCount(id)
}
