package api

import (
	"net/http"
	"strconv"
	"time"

	"noneland/backend/interview/internal/db"
	"noneland/backend/interview/internal/entity"
	"noneland/backend/interview/internal/pkg"

	"github.com/gin-gonic/gin"
)

// API 回應結構
type APIResponse struct {
	SpotBalance      *entity.BalanceResponse `json:"spot_balance"`
	ContractBalance  *entity.BalanceResponse `json:"contract_balance"`
}

// API 結構
type API struct {
	APIClient pkg.APIClient
	Cache     *pkg.Cache
	DBClient  db.DBClient
}

// NewAPI 創建新的 API
func NewAPI(apiClient pkg.APIClient, cache *pkg.Cache, dbClient db.DBClient) *API {
	return &API{
		APIClient: apiClient,
		Cache:     cache,
		DBClient:  dbClient,
	}
}

// GetBalances 處理函數
func (api *API) GetBalances(c *gin.Context) {
	cacheKey := "balances"

	if cachedData, found := api.Cache.Get(cacheKey); found {
		c.JSON(http.StatusOK, cachedData)
		return
	}

	spotBalance, err := api.APIClient.GetSpotBalance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contractBalance, err := api.APIClient.GetContractBalance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := APIResponse{
		SpotBalance:     spotBalance,
		ContractBalance: contractBalance,
	}

	api.Cache.Set(cacheKey, response, 5*time.Minute) // 緩存 5 分鐘

	c.JSON(http.StatusOK, response)
}

// FetchAndSaveSpotTransferRecords 獲取第三方交易所交易紀錄並保存到資料庫
func (api *API) FetchAndSaveSpotTransferRecords(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	current := c.DefaultQuery("current", "1")
	size := c.DefaultQuery("size", "10")

	startTimeInt, _ := strconv.ParseInt(startTime, 10, 64)
	endTimeInt, _ := strconv.ParseInt(endTime, 10, 64)
	currentInt, _ := strconv.Atoi(current)
	sizeInt, _ := strconv.Atoi(size)

	records, err := api.APIClient.GetSpotTransferRecords(startTimeInt, endTimeInt, currentInt, sizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 儲存交易紀錄
	err = api.DBClient.SaveTransactions(records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transactions have been saved to the database."})
}

// GetTransactions 從資料庫中獲取交易紀錄
func (api *API) GetTransactions(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	var startTimeInt, endTimeInt int64
	var err error

	if startTime != "" {
		startTimeInt, err = strconv.ParseInt(startTime, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startTime format"})
			return
		}
	}

	if endTime != "" {
		endTimeInt, err = strconv.ParseInt(endTime, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endTime format"})
			return
		}
	}

	transactions, err := api.DBClient.GetTransactions(startTimeInt, endTimeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}