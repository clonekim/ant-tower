//go:generate goversioninfo
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"twn-monitor/config"
	"twn-monitor/data"
	"twn-monitor/logger"
	"twn-monitor/server"
	"twn-monitor/sysagent"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()

	logger.Setup(cfg)

	if !cfg.UseConsole {
		sysagent.HideConsole()
		log.Info().Msg("Background mode: Console hidden")
	} else {
		log.Info().Msg("Console mode started")
	}

	if err := data.InitDB(cfg); err != nil {
		log.Fatal().Err(err).Msg("FATAL: failed to initialize database")
	}

	// Channel for WebSocket broadcast
	broadcast := make(chan interface{})

	// Initialize Monitor Service
	monService := sysagent.NewMonitorService(broadcast)

	// Create Context for Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Start Background Monitoring (with Context)
	go monService.Start(ctx)
	wsHub := server.NewHub(broadcast)
	go wsHub.Run()

	log.Info().Msg("Monitoring agent started")

	handler := server.NewHandler(wsHub)
	router := server.SetupRouter(handler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Server run in goroutine to allow signal handling
	go func() {
		log.Info().Str("port", cfg.Port).Msg("TWN server starting...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Server crashed")
		}
	}()

	// Wait for Interrupt Signal (Graceful Shutdown Trigger)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block here until signal received

	log.Info().Msg("Shutting down server...")

	// Cleanup
	cancel() // This stops the MonitorService loop

	// Shutdown HTTP Server
	ctxServer, cancelServer := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelServer()
	if err := srv.Shutdown(ctxServer); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited successfully")
}
