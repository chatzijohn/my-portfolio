package grpc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/chatzijohn/portfolio/apps/api/config"

	// Import generated proto
	pb "github.com/chatzijohn/portfolio/apps/api/internal/proto"
)

// New creates and configures a ready-to-start gRPC Server.
// Note: We don't pass 'router' or 'http.Handler' because gRPC handles routing internally via Protobufs.
func New(
	ctx context.Context,
	cfg *config.ServerConfig,
	pool *pgxpool.Pool,
) *googlegrpc.Server {

	// 1. Define Server Options (Timeouts & Keepalives)
	opts := []googlegrpc.ServerOption{
		googlegrpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 120 * time.Second, // IdleTimeout equivalent
			MaxConnectionAge:  5 * time.Minute,   // Force reconnect occasionally
			Time:              20 * time.Second,  // Ping client if idle
			Timeout:           5 * time.Second,   // Wait for ping ack
		}),
	}

	// 2. Create the raw gRPC Server
	grpcServer := googlegrpc.NewServer(opts...)

	// 3. Create YOUR implementation (from server.go)
	// We pass the dependencies your logic needs (DB, Prefs, etc)
	serviceImplementation := NewServer(pool)

	// 4. Register the implementation to the Server
	pb.RegisterPortfolioServiceServer(grpcServer, serviceImplementation)

	// 5. Enable Reflection (Crucial for debugging tools like Postman/gRPCurl)
	reflection.Register(grpcServer)

	return grpcServer
}
