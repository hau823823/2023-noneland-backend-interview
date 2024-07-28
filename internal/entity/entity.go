package entity

import "time"

// 定義餘額結構
type BalanceResponse struct {
	Asset   string  `json:"asset"`
	Balance float64 `json:"balance"`
}

// 定義交易紀錄結構
type Transaction struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}