package cache

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCache) Set(key string, value []byte, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}
