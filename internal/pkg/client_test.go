package pkg

import (
	"noneland/backend/interview/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMockAPIClient(t *testing.T) {
	mockAPIClient := new(MockAPIClient)

	mockSpotBalance := &entity.BalanceResponse{
		Asset:   "USDT",
		Balance: 1000.0,
	}
	mockContractBalance := &entity.BalanceResponse{
		Asset:   "USDT",
		Balance: 500.0,
	}
	mockTransactions := []entity.Transaction{
		{
			ID:        "1",
			Type:      "deposit",
			Amount:    100.0,
			Timestamp: time.Now(),
		},
	}

	// 設置 Mock 行為
	mockAPIClient.On("GetSpotBalance").Return(mockSpotBalance, nil)
	mockAPIClient.On("GetContractBalance").Return(mockContractBalance, nil)
	mockAPIClient.On("GetSpotTransactions").Return(mockTransactions, nil)

	spotBalance, err := mockAPIClient.GetSpotBalance()
	assert.NoError(t, err)
	assert.Equal(t, mockSpotBalance, spotBalance)

	contractBalance, err := mockAPIClient.GetContractBalance()
	assert.NoError(t, err)
	assert.Equal(t, mockContractBalance, contractBalance)

	transactions, err := mockAPIClient.GetSpotTransactions()
	assert.NoError(t, err)
	assert.Equal(t, mockTransactions, transactions)

	mockAPIClient.AssertExpectations(t)
}
