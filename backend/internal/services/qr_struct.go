package services

import (
	"database/sql"
	"errors"

	"github.com/phucnguyen/qrify/internal/models"
)

type QRCodeStore interface {
	Save(qr *models.QRCode) error
	FindByID(id string) (*models.QRCode, error)
	DeleteByID(id string) error
	FindByURL(url string) (*models.QRCode, error)
	IncrementScanCount(id string) error
}

type PostgresQRCodeStore struct {
	db *sql.DB
}

func NewPostgresQRCodeStore(db *sql.DB) *PostgresQRCodeStore {
	return &PostgresQRCodeStore{db: db}
}

func (s *PostgresQRCodeStore) Save(qr *models.QRCode) error {
	_, err := s.db.Exec(
		`INSERT INTO qr_codes (id, url, created_at, expires_at, image_base64, scan_count) VALUES ($1, $2, $3, $4, $5, $6)`,
		qr.ID, qr.URL, qr.CreatedAt, qr.ExpiresAt, qr.ImageBase64, qr.ScanCount,
	)
	return err
}

func (s *PostgresQRCodeStore) FindByID(id string) (*models.QRCode, error) {
	row := s.db.QueryRow(`SELECT id, url, created_at, expires_at, image_base64, scan_count FROM qr_codes WHERE id = $1`, id)
	var qr models.QRCode
	if err := row.Scan(&qr.ID, &qr.URL, &qr.CreatedAt, &qr.ExpiresAt, &qr.ImageBase64, &qr.ScanCount); err != nil {
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

func (s *PostgresQRCodeStore) FindByURL(url string) (*models.QRCode, error) {
	row := s.db.QueryRow(`SELECT id, url, created_at, expires_at, image_base64, scan_count FROM qr_codes WHERE url = $1`, url)
	var qr models.QRCode
	if err := row.Scan(&qr.ID, &qr.URL, &qr.CreatedAt, &qr.ExpiresAt, &qr.ImageBase64, &qr.ScanCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &qr, nil
}

func (s *PostgresQRCodeStore) IncrementScanCount(id string) error {
	_, err := s.db.Exec(`UPDATE qr_codes SET scan_count = scan_count + 1 WHERE id = $1`, id)
	return err
}
