package domain

type Pair struct {
	Pair            string // Pair key, lowercase pair name (e.g. usdtrub)
	AskUnit         string // Primary crypto
	BidUnit         string // Corresponding crypto
	MinAsk          int64  // Minimal ask amount (int with volume precision)
	MinBid          int64  // Minimal bid amount (int with volume precision)
	MakerFee        int64  // Maker fee (int with price precision)
	TakerFee        int64  // Taker fee (int with price precision)
	PricePrecision  int    // Precision for price values processing
	VolumePrecision int    // Precision for amounts value processing
	CreatedAt       int64  // UTC unixtime
	UpdatedAt       int64  // UTC unixtime
}
