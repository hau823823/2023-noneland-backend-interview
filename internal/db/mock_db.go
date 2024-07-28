package db

import (
	"noneland/backend/interview/internal/entity"

	"github.com/stretchr/testify/mock"
)

// MockDBClient 定義了一個模擬的資料庫客戶端
type MockDBClient struct {
	mock.Mock
}

func (m *MockDBClient) SaveTransactions(transactions []entity.Transaction) error {
	args := m.Called(transactions)
	return args.Error(0)
}
