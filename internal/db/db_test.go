package db

import (
	"noneland/backend/interview/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveTransactions(t *testing.T) {
	// 創建 Mock 資料庫客戶端
	mockDBClient := new(MockDBClient)

	mockTransactions := []entity.Transaction{
		{
			ID:        "1",
			Type:      "deposit",
			Amount:    100.0,
			Timestamp: time.Now(),
		},
	}

	// 設置 Mock 行為
	mockDBClient.On("SaveTransactions", mockTransactions).Return(nil)

	err := mockDBClient.SaveTransactions(mockTransactions)
	assert.NoError(t, err)
	mockDBClient.AssertExpectations(t)
}
