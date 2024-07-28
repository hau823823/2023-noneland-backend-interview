package pkg

import (
	"encoding/json"
	"noneland/backend/interview/internal/entity"
	"net/http"
)

// 定義 API 客戶端接口
type APIClient interface {
	GetSpotBalance() (*entity.BalanceResponse, error)
	GetContractBalance() (*entity.BalanceResponse, error)
	GetSpotTransactions() ([]entity.Transaction, error)
}

// 實現真實的 API 客戶端
type RealAPIClient struct{}

const (
	SpotBalanceEndpoint            = "https://api.xxexchange.com/spot/balance"
	ContractBalanceEndpoint        = "https://api.xxexchange.com/futures/balance"
	SpotTransactionHistoryEndpoint = "https://api.xxexchange.com/spot/transfer/records"
)

func (c *RealAPIClient) GetSpotBalance() (*entity.BalanceResponse, error) {
	resp, err := http.Get(SpotBalanceEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var balance entity.BalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, err
	}
	return &balance, nil
}

func (c *RealAPIClient) GetContractBalance() (*entity.BalanceResponse, error) {
	resp, err := http.Get(ContractBalanceEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var balance entity.BalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, err
	}
	return &balance, nil
}

func (c *RealAPIClient) GetSpotTransactions() ([]entity.Transaction, error) {
	resp, err := http.Get(SpotTransactionHistoryEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transactions []entity.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func NewRealAPIClient() *RealAPIClient {
	return &RealAPIClient{}
}
