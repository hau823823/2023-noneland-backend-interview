package api

import (
	"net/http"
	"noneland/backend/interview/internal/db"
	"noneland/backend/interview/internal/entity"
	"noneland/backend/interview/internal/pkg"
	"time"

	"github.com/gin-gonic/gin"
)

// API 回應結構
type APIResponse struct {
	SpotBalance      *entity.BalanceResponse `json:"spot_balance"`
	ContractBalance  *entity.BalanceResponse `json:"contract_balance"`
	SpotTransactions []entity.Transaction    `json:"spot_transactions"`
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

// GetBalancesAndTransactions 處理函數
func (api *API) GetBalancesAndTransactions(c *gin.Context) {
	cacheKey := "balances_and_transactions"

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

	spotTransactions, err := api.APIClient.GetSpotTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 儲存交易紀錄
	err = api.DBClient.SaveTransactions(spotTransactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := APIResponse{
		SpotBalance:      spotBalance,
		ContractBalance:  contractBalance,
		SpotTransactions: spotTransactions,
	}

	api.Cache.Set(cacheKey, response, 5*time.Minute) // 緩存 5 分鐘

	c.JSON(http.StatusOK, response)
}
