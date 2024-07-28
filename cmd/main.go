package main

import (
	"fmt"
	"log"
	"noneland/backend/interview/configs"
	"noneland/backend/interview/internal/api"
	"noneland/backend/interview/internal/db"
	"noneland/backend/interview/internal/pkg"
)

func main() {
	// 加載配置
	cfg := configs.NewConfig()

	// 初始化各個依賴
	apiClient := pkg.NewRealAPIClient()
	cache := pkg.NewCache()
	dbClient, err := db.NewMySQLClient(cfg.DSN)
	if err != nil {
		log.Fatalf("could not initialize db client: %v", err)
	}

	// 初始化 API
	apiInstance := api.NewAPI(apiClient, cache, dbClient)
	handlers := api.HandlerGroup{
		API: apiInstance,
	}

	// 啟動服務器
	server := api.NewServer(cfg, handlers)
	fmt.Printf("開始監聽 %v\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
