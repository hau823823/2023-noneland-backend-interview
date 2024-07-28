package pkg

import (
	"encoding/json"
	"errors"
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
type ExAPIClient struct{
	BaseURL string
	Client  *http.Client
}

func NewExAPIClient() *ExAPIClient {
	return &ExAPIClient{
		BaseURL: "https://api.example.com",
		Client:  &http.Client{},
	}
}

func (c *ExAPIClient) GetSpotBalance() (*entity.BalanceResponse, error) {
	url := c.BaseURL + "/spot/balance"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get spot balance")
	}

	var balance entity.BalanceResponse
	err = json.NewDecoder(resp.Body).Decode(&balance)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (c *ExAPIClient) GetContractBalance() (*entity.BalanceResponse, error) {
	url := c.BaseURL + "/futures/balance"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get contract balance")
	}

	var balance entity.BalanceResponse
	err = json.NewDecoder(resp.Body).Decode(&balance)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (c *ExAPIClient) GetSpotTransactions() ([]entity.Transaction, error) {
	url := c.BaseURL + "/spot/transfer/records"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get spot transactions")
	}

	var records struct {
		Rows []entity.Transaction `json:"rows"`
	}
	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return nil, err
	}

	return records.Rows, nil
}