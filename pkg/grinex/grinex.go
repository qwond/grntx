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

func (g *Grinex) GetMarkets() ([]Market, error) {
	resp, err := g.client.Get(g.baseURL + MarketsURL)
	if err != nil {
		return nil, fmt.Errorf("cannot make request:%e", err)
	}

	defer resp.Body.Close()

	var markets []Market

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&markets)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response:%e", err)
	}

	return markets, nil
}

func (g *Grinex) GetRate(pair string) (*Rate, error) {
	var resp rateResponse

	r, err := g.client.Get(g.baseURL + RateURL + "?market=" + pair)
	if err != nil {
		return nil, fmt.Errorf("cannot make request:%e", err)
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response:%e", err)
	}

	rate := RateFromDTO(resp, pair)
	return &rate, nil
}
