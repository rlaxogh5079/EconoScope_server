package main

import (
	"fmt"
	"net/http"

	"github.com/rlaxogh5079/EconoScope/api"
	"github.com/rlaxogh5079/EconoScope/config"
	"github.com/rlaxogh5079/EconoScope/pkg/logger"
)

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	logger.InitLogger()

	r := api.SetUpRouter()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logger.Log.Infof("Starting server on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatalf("server error: %v", err)
	}
}
