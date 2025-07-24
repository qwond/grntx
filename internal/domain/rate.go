package domain

type Rate struct {
	Pair        string
	AskUnit     string
	BidUnit     string
	AskPrice    int64
	BidPrice    int64
	Timestamp   int64
	Precision   int
	RetrievedAt int64
}
