package server

import (
	"url-shorter/internal/models"

	"github.com/stretchr/testify/mock"
)

type LinkServiceMock struct {
	mock.Mock
}

func (m *LinkServiceMock) GetLink(short string) (*models.Link, error) {
	args := m.Called(short)
	return args.Get(0).(*models.Link), args.Error(1)
}

func (m *LinkServiceMock) RemoveLink(short string) error {
	args := m.Called(short)
	return args.Error(0)
}


func (m *LinkServiceMock) CreateLink(link models.Link) error {
	args := m.Called(link)
	return args.Error(0)
}


func (m *LinkServiceMock) GetQRCode(short string) ([]byte, error) {
	args := m.Called(short)
	return args.Get(0).([]byte), args.Error(1)
}
