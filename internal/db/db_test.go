package db

import (
	"testing"
	"time"

	"noneland/backend/interview/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveTransactions(t *testing.T) {
	mockDBClient := NewMockDBClient()
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
	mockDBClient.On("SaveTransactions", mock.Anything).Return(nil)

	// 呼叫 SaveTransactions 方法
	err := mockDBClient.SaveTransactions(mockTransactions)

	assert.NoError(t, err)
	mockDBClient.AssertExpectations(t)
}

func TestGetAllTransactions(t *testing.T) {
	mockDBClient := NewMockDBClient()
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
	mockDBClient.On("GetAllTransactions").Return(mockTransactions, nil)

	// 呼叫 GetAllTransactions 方法
	transactions, err := mockDBClient.GetAllTransactions()

	assert.NoError(t, err)
	assert.Equal(t, mockTransactions, transactions)

	mockDBClient.AssertExpectations(t)
}

func TestGetTransactions(t *testing.T) {
	mockDBClient := NewMockDBClient()
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

	startTime := mockTransactions[0].Timestamp.Unix()
	endTime := mockTransactions[1].Timestamp.Unix()

	// 設置 Mock 行為
	mockDBClient.On("GetTransactions", startTime, endTime).Return(mockTransactions, nil)

	// 呼叫 GetTransactions 方法
	transactions, err := mockDBClient.GetTransactions(startTime, endTime)

	assert.NoError(t, err)
	assert.Equal(t, mockTransactions, transactions)

	mockDBClient.AssertExpectations(t)
}
