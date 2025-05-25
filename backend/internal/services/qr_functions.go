package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/phucnguyen/qrify/internal/models"
	"github.com/skip2/go-qrcode"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type QRCodeStore interface {
	Save(qr *models.QRCode) error
	FindByID(id string) (*models.QRCode, error)
	DeleteByID(id string) error
}

type PostgresQRCodeStore struct {
	db *sql.DB
}

func NewPostgresQRCodeStore(db *sql.DB) *PostgresQRCodeStore {
	return &PostgresQRCodeStore{db: db}
}

func (s *PostgresQRCodeStore) Save(qr *models.QRCode) error {
	_, err := s.db.Exec(
		`INSERT INTO qr_codes (id, url, created_at, expires_at, image_base64) VALUES ($1, $2, $3, $4, $5)`,
		qr.ID, qr.URL, qr.CreatedAt, qr.ExpiresAt, qr.ImageBase64,
	)
	return err
}

func (s *PostgresQRCodeStore) FindByID(id string) (*models.QRCode, error) {
	row := s.db.QueryRow(`SELECT id, url, created_at, expires_at, image_base64 FROM qr_codes WHERE id = $1`, id)
	var qr models.QRCode
	if err := row.Scan(&qr.ID, &qr.URL, &qr.CreatedAt, &qr.ExpiresAt, &qr.ImageBase64); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &qr, nil
}

func (s *PostgresQRCodeStore) DeleteByID(id string) error {
	_, err := s.db.Exec(`DELETE FROM qr_codes WHERE id = $1`, id)
	return err
}

type QRService struct {
	store QRCodeStore
}

func NewQRService(store QRCodeStore) *QRService {
	return &QRService{
		store:        store,
	}
}

func addTextBelow(img image.Image, text string) (image.Image, error) {
	qrBounds := img.Bounds()
	textHeight := 20
	gap := 20       
	newHeight := qrBounds.Dy() + gap + textHeight
	newImg := image.NewRGBA(image.Rect(0, 0, qrBounds.Dx(), newHeight))

	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, qrBounds, img, image.Point{}, draw.Over)

	
	col := color.Black
	point := fixed.Point26_6{
		X: fixed.I((qrBounds.Dx() - len(text)*7) / 2), 
		Y: fixed.I(qrBounds.Dy() + gap + 15),          
	}
	d := &font.Drawer{
		Dst:  newImg,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
	return newImg, nil
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
		expiresAt = time.Now().Add(time.Duration(req.ExpiresInSec) * time.Second)
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
		ID:        qr.ID,
		URL:       qr.URL,
		QRCodeURL: "/r/" + qr.ID,
		CreatedAt: qr.CreatedAt,
		ExpiresAt: qr.ExpiresAt,
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

	if !qr.ExpiresAt.IsZero() && time.Now().After(qr.ExpiresAt) {
		return nil, errors.New("QR code has expired")
	}
	return &models.QRCodeResponse{
		ID:        qr.ID,
		URL:       qr.URL,
		QRCodeURL: "/r/" + qr.ID,
		CreatedAt: qr.CreatedAt,
		ExpiresAt: qr.ExpiresAt,
		
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