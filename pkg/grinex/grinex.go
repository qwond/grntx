// Grinex client package
package grinex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

func (g *Grinex) GetMarkets() ([]MarketDTO, error) {
	resp, err := g.client.Get(g.baseURL + MarketsURL)
	if err != nil {
		return nil, fmt.Errorf("cannot make request:%e", err)
	}

	defer resp.Body.Close()

	var markets []MarketDTO

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&markets)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response:%e", err)
	}

	return markets, nil
}

func (g *Grinex) GetRate(currency string) (RateDTO, error) {
	var rate RateDTO

	resp, err := g.client.Get(g.baseURL + RateURL + "?market=" + currency)
	if err != nil {
		return rate, fmt.Errorf("cannot make request:%e", err)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&rate)
	if err != nil {
		return rate, fmt.Errorf("cannot decode response:%e", err)
	}

	return rate, nil
}
