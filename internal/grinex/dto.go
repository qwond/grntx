package grinex

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/qwond/grntx/internal/domain"
)

type pairResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AskUnit        string `json:"ask_unit"`
	BidUnit        string `json:"bid_unit"`
	PricePrecision int    `json:"price_precision"`
}

func PairFromDTO(p pairResponse) domain.Pair {
	return domain.Pair{
		Pair:           p.ID,
		AskUnit:        p.AskUnit,
		BidUnit:        p.BidUnit,
		PricePrecision: p.PricePrecision,
	}
}

// Rates response is a bit more complex than required
// so request object needs addition conversion
type rateResponse struct {
	Timestamp int64 `json:"timestamp"`
	Asks      []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
		Amount string `json:"amount"`
		Factor string `json:"factor"`
		Type   string `json:"type"`
	} `json:"asks"`
	Bids []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
		Amount string `json:"amount"`
		Factor string `json:"factor"`
		Type   string `json:"type"`
	} `json:"bids"`
}

// RateFromDTO converts rateResponse into domain entity
func RateFromDTO(r rateResponse, pair domain.Pair) (domain.Rate, error) {
	var (
		rate domain.Rate
		err  error
	)

	// Check if we have enough rates to convert
	if len(r.Asks) < 1 || len(r.Bids) < 1 {
		return rate, fmt.Errorf("missing required info")
	}

	rate = domain.Rate{
		Pair:        pair.Pair,
		AskUnit:     pair.AskUnit,
		BidUnit:     pair.BidUnit,
		Timestamp:   r.Timestamp,
		RetrievedAt: time.Now().Unix(),
	}

	rate.AskPrice, err = StrFloatToInt(r.Asks[0].Price, pair.PricePrecision)
	if err != nil {
		return rate, fmt.Errorf("cannot convert response to Rate:%e", err)
	}

	rate.BidPrice, err = StrFloatToInt(r.Bids[0].Price, pair.PricePrecision)
	if err != nil {
		return rate, fmt.Errorf("cannot convert response to Rate:%e", err)
	}

	return rate, nil
}

// StrFloatToInt converts string price, represented as float
// into in64 with according precision.
func StrFloatToInt(str string, prec int) (int64, error) {
	fl, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert(%s) into float:%e", str, err)
	}

	// multiply float by 10 power of preceision
	return int64(fl * math.Pow(10, float64(prec))), nil
}
