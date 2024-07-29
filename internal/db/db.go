package db

import (
	"time"

	"noneland/backend/interview/internal/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBClient 定義資料庫客戶端接口
type DBClient interface {
	SaveTransactions(transactions []entity.Transaction) error
	GetAllTransactions() ([]entity.Transaction, error)
	GetTransactions(startTime, endTime int64) ([]entity.Transaction, error)
}

// MySQLClient 實現 DBClient 接口
type MySQLClient struct {
	DB *gorm.DB
}


func NewMySQLDBClient(dsn string) (DBClient, error) {
	return NewMySQLClient(dsn)
}

// NewMySQLClient 創建一個新的 MySQLClient
func NewMySQLClient(dsn string) (*MySQLClient, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entity.Transaction{})
	if err != nil {
		return nil, err
	}

	return &MySQLClient{DB: db}, nil
}

// SaveTransactions 儲存交易紀錄
func (c *MySQLClient) SaveTransactions(transactions []entity.Transaction) error {
	for _, txn := range transactions {
		err := c.DB.Create(&txn).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllTransactions 獲取所有交易紀錄
func (c *MySQLClient) GetAllTransactions() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := c.DB.Find(&transactions).Error
	return transactions, err
}

func (c *MySQLClient) GetTransactions(startTime, endTime int64) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	query := c.DB

	if startTime > 0 {
		query = query.Where("timestamp >= ?", time.Unix(startTime, 0))
	}
	if endTime > 0 {
		query = query.Where("timestamp <= ?", time.Unix(endTime, 0))
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
