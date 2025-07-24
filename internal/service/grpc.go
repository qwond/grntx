package service

import (
	"context"
	"time"

	v1 "github.com/qwond/grntx/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRates implements v1.RatesServiceServer.
func (rs *RatesService) GetRates(
	ctx context.Context,
	req *v1.GetRatesRequest,
) (*v1.GetRatesResponse, error) {
	symbol := req.GetSymbol()
	if symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "pair is required")
	}

	rs.log.Info("request", zap.String("symbol", symbol))
	// Enumerate all available pairs and collect pairs
	// with Ask unit == symbol

	rates := []*v1.Rate{}

	for _, pair := range rs.pairs {
		if pair.AskUnit == symbol {
			// 1. Check cached record
			rate, ok := rs.rates[pair.Pair]
			if !ok || time.Now().Unix()-rate.RetrievedAt > 10 {
				// not in cache, retrieve
				rs.log.Info("not cached", zap.String("symbol", symbol))

				rate, err := rs.grinex.GetRate(pair)
				if err != nil {
					// cannot retrieve, bail out
					return nil, status.Error(codes.Internal, "retrieving failed, check logs")
				}
				rs.rates[pair.Pair] = *rate
			}

			rates = append(rates, &v1.Rate{
				Pair:      rate.Pair,
				AskUnit:   rate.AskUnit,
				BidUnit:   rate.BidUnit,
				AskPrice:  rate.AskPrice,
				BidPrice:  rate.BidPrice,
				Timestamp: rate.Timestamp,
				Precision: int64(rate.Precision),
			})
		}
	}
	resp := v1.GetRatesResponse{
		Rates: rates,
	}

	return &resp, nil
}

// HealthCheck implements v1.RatesServiceServer.
func (r *RatesService) HealthCheck(ctx context.Context, req *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{Status: "SERVING"}, nil
}
