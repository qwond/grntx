package domain

import "time"

type Rate struct {
	Pair      string
	Bid       float64
	Ask       float64
	Timestamp time.Time
}
