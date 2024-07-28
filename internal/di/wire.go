//go:build wireinject
// +build wireinject

package di

import (
	"noneland/backend/interview/configs"
	"noneland/backend/interview/internal/api"
	"noneland/backend/interview/internal/db"
	"noneland/backend/interview/internal/pkg"
	"sync"

	"github.com/google/wire"
)

var (
	configOnce sync.Once
	cfg        *configs.Config

	dbOnce   sync.Once
	database *db.MySQLClient
)

// NewConfig 確保配置僅初始化一次
func NewConfig() *configs.Config {
	configOnce.Do(func() {
		cfg = configs.NewConfig()
	})
	return cfg
}

// NewDBClient 確保數據庫連接僅初始化一次
func NewDBClient(cfg *configs.Config) (db.DBClient, error) {
	var err error
	dbOnce.Do(func() {
		database, err = db.NewMySQLClient(cfg.DSN)
	})
	return database, err
}

// NewExAPIClient 構造真實 API 客戶端
func NewExAPIClient() pkg.APIClient {
	return pkg.NewExAPIClient()
}

// NewAPI 構造真實的 API 實例
func NewAPI() (*api.API, error) {
	wire.Build(api.NewAPI, NewExAPIClient, pkg.NewCache, NewDBClient, NewConfig)
	return &api.API{}, nil
}

// NewMockDBClient 構造 Mock DB 客戶端
/**
func NewMockDBClient() db.DBClient {
	return db.NewMockDBClient()
}

// NewMockAPIClient 構造 Mock API 客戶端
func NewMockAPIClient() pkg.APIClient {
	return pkg.NewMockAPIClient()
}

// NewMockAPI 構造 Mock API 實例
func NewMockAPI() (*api.API, error) {
	wire.Build(api.NewAPI, NewMockAPIClient, pkg.NewCache, NewMockDBClient, NewConfig)
	return &api.API{}, nil
}
*/
