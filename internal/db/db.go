package db

import (
	"noneland/backend/interview/internal/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBClient 定義資料庫客戶端接口
type DBClient interface {
	SaveTransactions(transactions []entity.Transaction) error
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