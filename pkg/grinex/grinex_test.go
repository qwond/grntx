package grinex_test

import (
	"testing"

	"github.com/qwond/grntx/pkg/grinex"
)

// TODO mock server
func TestGetMarkets(t *testing.T) {
	g := grinex.New("https://grinex.io")
	_, err := g.GetMarkets()
	if err != nil {
		t.Errorf("cannot retrieve market:%e", err)
	}
}

// TODO mock server
func TestGetRates(t *testing.T) {
	g := grinex.New("https://grinex.io")
	_, err := g.GetRate("usdtrub")
	if err != nil {
		t.Errorf("cannot retrieve rates:%e", err)
	}
}
