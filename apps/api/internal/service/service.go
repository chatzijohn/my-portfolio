package service

import (
	"github.com/chatzijohn/portfolio/apps/api/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	WaterSupply() WaterSupplyService
}

type service struct {
	waterSupply WaterSupplyService
}

func (s *service) WaterSupply() WaterSupplyService {
	return s.waterSupply
}

func New(pool *pgxpool.Pool) Service {
	store := repository.New(pool)
	return &service{
		waterSupply: NewWaterSupplyService(store.WaterSupply()),
	}
}
