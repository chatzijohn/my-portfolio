package grpc

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	// Import generated proto
	pb "github.com/chatzijohn/portfolio/apps/api/internal/proto"
)

// Server implements the gRPC interface
type Server struct {
	pb.UnimplementedPortfolioServiceServer
	DB *pgxpool.Pool
}

// NewServer creates the implementation (Logic)
// This is the function the compiler was missing!
func NewServer(db *pgxpool.Pool) *Server {
	return &Server{
		DB: db,
	}
}

// GetHero implements the PortfolioService.GetHero RPC
func (s *Server) GetHero(ctx context.Context, req *pb.Empty) (*pb.HeroResponse, error) {
	// Mock response for now
	return &pb.HeroResponse{
		Headline:           "Hi, I'm Panagiotis",
		Subheadline:        "Software & DevOps Engineer",
		IsAvailableForWork: true,
	}, nil
}

func (s *Server) GetPosts(ctx context.Context, req *pb.Empty) (*pb.PostList, error) {
	return &pb.PostList{
		Posts: []*pb.BlogPost{},
	}, nil
}
