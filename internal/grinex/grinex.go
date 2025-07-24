// Grinex client package
package grinex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/qwond/grntx/internal/domain"
)

const (
	MarketsURL = "/api/v2/markets"
	RateURL    = "/api/v2/depth"
)

type Grinex struct {
	baseURL string
	client  http.Client
}

func New(baseURL string) *Grinex {
	return &Grinex{
		baseURL: baseURL,
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *Grinex) GetMarkets() ([]domain.Pair, error) {
	resp, err := g.client.Get(g.baseURL + MarketsURL)
	if err != nil {
		return nil, fmt.Errorf("cannot make request:%e", err)
	}

	//nolint: all
	defer resp.Body.Close()

	var data []pairResponse

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response:%e", err)
	}

	pairs := []domain.Pair{}
	for _, item := range data {
		pairs = append(pairs, PairFromDTO(item))
	}

	return pairs, nil
}

func (g *Grinex) GetRate(pair domain.Pair) (*domain.Rate, error) {
	var resp rateResponse

	r, err := g.client.Get(g.baseURL + RateURL + "?market=" + pair.Pair)
	if err != nil {
		return nil, fmt.Errorf("cannot make request:%e", err)
	}
	// nolint: all
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response:%e", err)
	}

	rate, err := RateFromDTO(resp, pair)
	if err != nil {
		return nil, err
	}

	return &rate, nil
}
