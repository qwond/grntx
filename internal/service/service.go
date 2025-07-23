package service

import (
	"context"

	v1 "github.com/qwond/grntx/api/v1"
	"github.com/qwond/grntx/pkg/grinex"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RatesService struct {
	v1.UnimplementedRatesServiceServer
	grinex *grinex.Grinex
}

// GetRates implements v1.RatesServiceServer.
func (rs *RatesService) GetRates(ctx context.Context, req *v1.GetRatesRequest) (*v1.GetRatesResponse, error) {
	pair := req.GetPair()
	if pair == "" {
		return nil, status.Error(codes.InvalidArgument, "pair is required")
	}

	rate, err := rs.grinex.GetRate(pair)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot retrieve rate for pair")
	}

	return &v1.GetRatesResponse{
		Pair:      pair,
		Ask:       rate.AskPrice,
		Bid:       rate.BidPrice,
		Timestamp: int64(rate.Timestamp),
	}, nil
}

func New(grnx *grinex.Grinex) *RatesService {
	return &RatesService{
		grinex: grnx,
	}
}

// HealthCheck implements v1.RatesServiceServer.
func (r *RatesService) HealthCheck(ctx context.Context, req *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{Status: "SERVING"}, nil
}
