package pkg

import (
	"noneland/backend/interview/internal/entity"

	"github.com/stretchr/testify/mock"
)

// MockAPIClient 定義了一個模擬的 API 客戶端
type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) GetSpotBalance() (*entity.BalanceResponse, error) {
	args := m.Called()
	return args.Get(0).(*entity.BalanceResponse), args.Error(1)
}

func (m *MockAPIClient) GetContractBalance() (*entity.BalanceResponse, error) {
	args := m.Called()
	return args.Get(0).(*entity.BalanceResponse), args.Error(1)
}

func (m *MockAPIClient) GetSpotTransferRecords(startTime, endTime int64, current, size int) ([]entity.Transaction, error) {
	args := m.Called(startTime, endTime, current, size)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func NewMockAPIClient() *MockAPIClient {
	return &MockAPIClient{}
}