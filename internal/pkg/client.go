package pkg

import (
	"encoding/json"
	"net/http"
	"noneland/backend/interview/internal/entity"
)

// 定義 API 客戶端接口
type APIClient interface {
	GetSpotBalance() (*entity.BalanceResponse, error)
	GetContractBalance() (*entity.BalanceResponse, error)
	GetSpotTransactions() ([]entity.Transaction, error)
}

// 實現第三方交易所的 API 客戶端
type ExAPIClient struct{}

const (
	SpotBalanceEndpoint            = "https://api.xxexchange.com/spot/balance"
	ContractBalanceEndpoint        = "https://api.xxexchange.com/futures/balance"
	SpotTransactionHistoryEndpoint = "https://api.xxexchange.com/spot/transfer/records"
)

func NewExAPIClient() *ExAPIClient {
	return &ExAPIClient{}
}

func (c *ExAPIClient) GetSpotBalance() (*entity.BalanceResponse, error) {
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

func (c *ExAPIClient) GetContractBalance() (*entity.BalanceResponse, error) {
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

func (c *ExAPIClient) GetSpotTransactions() ([]entity.Transaction, error) {
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
