package grinex

type MarketDTO struct {
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

type RateDTO struct {
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
