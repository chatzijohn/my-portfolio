package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chatzijohn/portfolio/apps/api/config"
	"github.com/chatzijohn/portfolio/apps/api/internal/db"
	grpcServer "github.com/chatzijohn/portfolio/apps/api/internal/transport/grpc"
)

func main() {
	// 1. Setup Structured Logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	// 2. Create root context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 3. Database Connection
	pool, err := db.NewPGPool(ctx, &cfg.DB)
	if err != nil {
		logger.Error("Unable to connect to DB", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("Connected to Database", "pool_config", cfg.DB)

	// 4. Create TCP Listener (REQUIRED for gRPC)
	port := cfg.SERVER.PORT
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	// 5. Initialize gRPC Server
	srv := grpcServer.New(ctx, &cfg.SERVER, pool)

	// 6. Run Server in Goroutine
	go func() {
		logger.Info("gRPC Server started", "port", port)
		if err := srv.Serve(lis); err != nil {
			logger.Error("gRPC server failed", "error", err)
			stop()
		}
	}()

	// 7. Wait for Shutdown Signal
	<-ctx.Done()
	logger.Info("Shutdown signal received")

	// 8. Graceful Shutdown
	// Create a channel to wait for GracefulStop to finish
	stopped := make(chan struct{})
	go func() {
		srv.GracefulStop()
		close(stopped)
	}()

	// Force stop after 5 seconds if GracefulStop hangs
	select {
	case <-stopped:
		logger.Info("Server stopped gracefully")
	case <-time.After(5 * time.Second):
		logger.Error("Shutdown timed out, forcing stop")
		srv.Stop()
	}
}
