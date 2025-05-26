package tests

import (
	"errors"

	"github.com/phucnguyen/qrify/internal/models"
)

// MockQRCodeStore implements services.QRCodeStore
type MockQRCodeStore struct {
	qrCodes map[string]*models.QRCode
}

func NewMockQRCodeStore() *MockQRCodeStore {
	return &MockQRCodeStore{
		qrCodes: make(map[string]*models.QRCode),
	}
}

func (m *MockQRCodeStore) Save(qr *models.QRCode) error {
	m.qrCodes[qr.ID] = qr
	return nil
}

func (m *MockQRCodeStore) FindByID(id string) (*models.QRCode, error) {
	qr, ok := m.qrCodes[id]
	if !ok {
		return nil, errors.New("QR code not found")
	}
	return qr, nil
}

func (m *MockQRCodeStore) DeleteByID(id string) error {
	delete(m.qrCodes, id)
	return nil
}

func (m *MockQRCodeStore) FindByURL(url string) (*models.QRCode, error) {
	for _, qr := range m.qrCodes {
		if qr.URL == url {
			return qr, nil
		}
	}
	return nil, nil
}

func (m *MockQRCodeStore) IncrementScanCount(id string) error {
	qr, ok := m.qrCodes[id]
	if !ok {
		return errors.New("QR code not found")
	}
	qr.ScanCount++
	return nil
}
