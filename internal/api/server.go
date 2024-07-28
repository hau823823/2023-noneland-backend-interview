package api

import (
	"fmt"
	"net/http"
	"time"

	"noneland/backend/interview/configs"

	"github.com/gin-gonic/gin"
)

type HandlerGroup struct {
	*API
}

func registerRoute(router *gin.Engine, handlers HandlerGroup) {
	v1 := router.Group("/api/v1")

	v1.GET("/balances_and_transactions", handlers.API.GetBalancesAndTransactions)
}

func setupServer(router *gin.Engine, cfg *configs.Config) *http.Server {
	if cfg.Mode != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Port),
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s
}

func NewServer(cfg *configs.Config, handlers HandlerGroup) *http.Server {
	router := gin.New()
	server := setupServer(router, cfg)
	registerRoute(router, handlers)
	return server
}

// NewRouter 可用在 httptest.NewServer 進行 integration test
func NewRouter(cfg *configs.Config, handlers HandlerGroup) *gin.Engine {
	router := gin.New()
	setupServer(router, cfg)
	registerRoute(router, handlers)
	return router
}
