package services

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

func (m *LinkServiceMock) RemoveLink(creatorLogin string, short string) error {
	args := m.Called(creatorLogin, short)
	return args.Error(0)
}

func (m *LinkServiceMock) CreateLink(link models.Link) error {
	args := m.Called(link)
	return args.Error(0)
}

func (m *LinkServiceMock) EditLink(creatorLogin string, short string, editedLink models.Link) error {
	args := m.Called(creatorLogin, short, editedLink)
	return args.Error(0)
}

func (m *LinkServiceMock) GetQRCode(link string, size int) ([]byte, error) {
	args := m.Called(link, size)
	return args.Get(0).([]byte), args.Error(1)
}
