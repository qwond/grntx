package service

import (
	"fmt"

	v1 "github.com/qwond/grntx/api/v1"
	"github.com/qwond/grntx/internal/domain"
	"github.com/qwond/grntx/internal/grinex"
	"github.com/qwond/grntx/internal/repository"
	"go.uber.org/zap"
)

type RatesService struct {
	log    *zap.Logger
	pairs  map[string]domain.Pair // available pairs
	rates  map[string]domain.Rate // rates cache
	repo   *repository.Repository
	grinex *grinex.Grinex

	v1.UnimplementedRatesServiceServer
}

// Creates new RateService instance.
func New(log *zap.Logger, grnx *grinex.Grinex, repo *repository.Repository) *RatesService {
	rs := &RatesService{
		log:    log,
		pairs:  make(map[string]domain.Pair),
		rates:  make(map[string]domain.Rate),
		repo:   repo,
		grinex: grnx,
	}
	rs.WarmUp()
	return rs
}

// WarmUp - RateService preparation:
// - Retrieving fresh pairs list from remote and store into db
// - Prepare caches
func (rs *RatesService) WarmUp() error {
	prs, err := rs.grinex.GetMarkets()
	if err != nil {
		return fmt.Errorf("cannot proceed without initial pairs list:%e", err)
	}

	for _, pr := range prs {
		rs.pairs[pr.Pair] = pr
	}

	return nil
}
