package pkg

import (
	"testing"
	"time"

	"noneland/backend/interview/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestGetSpotBalance(t *testing.T) {
	mockAPIClient := NewMockAPIClient()
	mockResponse := &entity.BalanceResponse{
		Asset:   "USDT",
		Balance: 1000.0,
	}

	// 設置 Mock 行為
	mockAPIClient.On("GetSpotBalance").Return(mockResponse, nil)

	// 呼叫 GetSpotBalance 方法
	balance, err := mockAPIClient.GetSpotBalance()

	assert.NoError(t, err)
	assert.Equal(t, mockResponse, balance)

	mockAPIClient.AssertExpectations(t)
}

func TestGetContractBalance(t *testing.T) {
	mockAPIClient := NewMockAPIClient()
	mockResponse := &entity.BalanceResponse{
		Asset:   "USDT",
		Balance: 500.0,
	}

	// 設置 Mock 行為
	mockAPIClient.On("GetContractBalance").Return(mockResponse, nil)

	// 呼叫 GetContractBalance 方法
	balance, err := mockAPIClient.GetContractBalance()

	assert.NoError(t, err)
	assert.Equal(t, mockResponse, balance)

	mockAPIClient.AssertExpectations(t)
}

func TestGetSpotTransferRecords(t *testing.T) {
	mockAPIClient := NewMockAPIClient()
	mockTransactions := []entity.Transaction{
		{
			ID:        "1",
			Type:      "IN",
			Amount:    0.1,
			Asset:     "BNB",
			Status:    "CONFIRMED",
			Timestamp: time.Unix(1566898617, 0),
			TxID:      5240372201,
		},
		{
			ID:        "2",
			Type:      "OUT",
			Amount:    5.0,
			Asset:     "USDT",
			Status:    "CONFIRMED",
			Timestamp: time.Unix(1566888436, 0),
			TxID:      5239810406,
		},
	}

	// 設置 Mock 行為
	mockAPIClient.On("GetSpotTransferRecords", int64(0), int64(0), 1, 10).Return(mockTransactions, nil)

	// 呼叫 GetSpotTransferRecords 方法
	records, err := mockAPIClient.GetSpotTransferRecords(0, 0, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, mockTransactions, records)

	mockAPIClient.AssertExpectations(t)
}
