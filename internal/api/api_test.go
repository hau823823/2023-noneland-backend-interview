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
	"github.com/stretchr/testify/mock"
)

func TestGetBalances(t *testing.T) {
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

	// 設置 Mock 行為
	mockAPIClient.On("GetSpotBalance").Return(mockSpotBalance, nil)
	mockAPIClient.On("GetContractBalance").Return(mockContractBalance, nil)

	apiInstance := NewAPI(mockAPIClient, mockCache, mockDBClient)

	// 創建 Gin 引擎並註冊路由
	router := gin.Default()
	router.GET("/api/v1/balances", apiInstance.GetBalances)

	// 發送 HTTP 請求並驗證響應
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/balances", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockSpotBalance, response.SpotBalance)
	assert.Equal(t, mockContractBalance, response.ContractBalance)

	// 驗證第二次請求應從緩存中獲取
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockAPIClient.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestFetchAndSaveSpotTransferRecords(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 創建 Mock 客戶端
	mockAPIClient := new(pkg.MockAPIClient)
	mockDBClient := new(db.MockDBClient)

	// 定義 Mock 返回值
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
	mockDBClient.On("SaveTransactions", mock.Anything).Return(nil)

	apiInstance := NewAPI(mockAPIClient, pkg.NewCache(), mockDBClient)

	// 創建 Gin 引擎並註冊路由
	router := gin.Default()
	router.POST("/api/v1/spot_transfer_records", apiInstance.FetchAndSaveSpotTransferRecords)

	// 發送 HTTP 請求並驗證響應
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/spot_transfer_records", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Transactions have been saved to the database.", response["message"])

	mockAPIClient.AssertExpectations(t)
	mockDBClient.AssertExpectations(t)
}

func TestGetTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 創建 Mock 客戶端
	mockDBClient := new(db.MockDBClient)

	// 定義 Mock 返回值
	mockTransactions := []entity.Transaction{
		{
			ID:        "1",
			Type:      "IN",
			Amount:    0.1,
			Asset:     "BNB",
			Status:    "CONFIRMED",
			Timestamp: time.Now().UTC(),
			TxID:      5240372201,
		},
		{
			ID:        "2",
			Type:      "OUT",
			Amount:    5.0,
			Asset:     "USDT",
			Status:    "CONFIRMED",
			Timestamp: time.Now().UTC(),
			TxID:      5239810406,
		},
	}

	// 設置 Mock 行為
	mockDBClient.On("GetTransactions", int64(0), int64(0)).Return(mockTransactions, nil)

	apiInstance := NewAPI(nil, nil, mockDBClient)

	// 創建 Gin 引擎並註冊路由
	router := gin.Default()
	router.GET("/api/v1/transactions", apiInstance.GetTransactions)

	// 發送 HTTP 請求並驗證響應
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/transactions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []entity.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockTransactions, response)

	mockDBClient.AssertExpectations(t)
}
