package grinex

type Market struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	AskUnit         string `json:"ask_unit"`
	BidUnit         string `json:"bid_unit"`
	MinAsk          string `json:"min_ask"`
	MinBid          string `json:"min_bid"`
	MakerFee        string `json:"maker_fee"`
	TakerFee        string `json:"taker_fee"`
	PricePrecision  int    `json:"price_precision"`
	VolumePrecision int    `json:"volume_precision"`
}

// Rates response is abit more complex than required
// so request object needs addition conversion
type rateResponse struct {
	Timestamp int `json:"timestamp"`
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

type Rate struct {
	Pair      string
	AskPrice  string
	BidPrice  string
	Timestamp int
}

func RateFromDTO(r rateResponse, pair string) Rate {
	return Rate{
		Pair:      pair,
		AskPrice:  r.Asks[0].Price,
		BidPrice:  r.Bids[0].Price,
		Timestamp: r.Timestamp,
	}

}
