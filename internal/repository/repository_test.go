package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/qwond/grntx/internal/domain"
	"github.com/qwond/grntx/internal/repository"
)

func setupRepo() (*repository.Repository, []domain.Pair, error) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return nil, nil, fmt.Errorf("please set DSN environment variable")
	}

	pairAAABBB := domain.Pair{
		Pair:            "aaabbb",
		AskUnit:         "AAA",
		BidUnit:         "BBB",
		MinAsk:          3,
		MinBid:          3,
		MakerFee:        1,
		TakerFee:        1,
		PricePrecision:  2,
		VolumePrecision: 2,
	}

	pairXXXYYY := domain.Pair{
		Pair:            "xxxyyyy",
		AskUnit:         "XXX",
		BidUnit:         "YYY",
		MinAsk:          3,
		MinBid:          3,
		MakerFee:        1,
		TakerFee:        1,
		PricePrecision:  2,
		VolumePrecision: 2,
	}

	pairs := []domain.Pair{pairAAABBB, pairXXXYYY}

	repo, err := repository.New(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create repository")
	}

	return repo, pairs, nil
}

func TestRepository(t *testing.T) {
	repo, pairs, err := setupRepo()
	if err != nil {
		t.Fatalf("cannot setup repository:%v", err)
	}

	ctx := context.Background()
	isNew, err := repo.PairUpsert(ctx, pairs[1])
	if err != nil {
		t.Errorf("cannot upsert pair:%v", err)
	}
	if !isNew {
		t.Error("first insert not returns new flag")
	}

	// sleep for 2 sec for check updated_at update
	isNew, err = repo.PairUpsert(ctx, pairs[1])
	if err != nil {
		t.Error("cannot update pair")
	}
	if isNew {
		t.Error("update returns new flag")
	}

	savedPairs, err := repo.PairsList(ctx)
	if err != nil {
		t.Errorf("cannot retrieve pairs list:%v", err)
	}
	if savedPairs[0].UpdatedAt == savedPairs[0].CreatedAt {
		t.Errorf("updated at is not updating")
	}
}
