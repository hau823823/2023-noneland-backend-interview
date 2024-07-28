package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"noneland/backend/interview/internal/db"
	"noneland/backend/interview/internal/entity"
	"noneland/backend/interview/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetBalancesAndTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 創建 Mock 客戶端
	mockAPIClient := new(pkg.MockAPIClient)
	mockDBClient := new(db.MockDBClient)
	mockCache := pkg.NewCache()

	// 定義 Mock 返回值
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
			Timestamp: time.Now().UTC(),
		},
	}

	// 設置 Mock 行為
	mockAPIClient.On("GetSpotBalance").Return(mockSpotBalance, nil)
	mockAPIClient.On("GetContractBalance").Return(mockContractBalance, nil)
	mockAPIClient.On("GetSpotTransactions").Return(mockTransactions, nil)
	mockDBClient.On("SaveTransactions", mockTransactions).Return(nil)

	apiInstance := NewAPI(mockAPIClient, mockCache, mockDBClient)

	// 創建 Gin 引擎並註冊路由
	router := gin.Default()
	router.GET("/api/v1/balances_and_transactions", apiInstance.GetBalancesAndTransactions)

	// 發送 HTTP 請求並驗證響應
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/balances_and_transactions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockSpotBalance, response.SpotBalance)
	assert.Equal(t, mockContractBalance, response.ContractBalance)
	assert.Equal(t, mockTransactions, response.SpotTransactions)

	// 驗證第二次請求應從緩存中獲取
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockAPIClient.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestGetAllTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 創建 Mock 客戶端
	mockDBClient := new(db.MockDBClient)

	// 定義 Mock 返回值
	mockTransactions := []entity.Transaction{
		{
			ID:        "1",
			Type:      "deposit",
			Amount:    100.0,
			Timestamp: time.Now().UTC(), // 統一使用 UTC 時區
		},
		{
			ID:        "2",
			Type:      "withdrawal",
			Amount:    50.0,
			Timestamp: time.Now().UTC(), // 統一使用 UTC 時區
		},
	}

	// 設置 Mock 行為
	mockDBClient.On("GetAllTransactions").Return(mockTransactions, nil)

	apiInstance := NewAPI(nil, nil, mockDBClient)

	// 創建 Gin 引擎並註冊路由
	router := gin.Default()
	router.GET("/api/v1/all_transactions", apiInstance.GetAllTransactions)

	// 發送 HTTP 請求並驗證響應
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/all_transactions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []entity.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockTransactions, response)

	mockDBClient.AssertExpectations(t)
}