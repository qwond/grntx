package grinex_test

import (
	"testing"

	"github.com/qwond/grntx/internal/grinex"
)

func TestStrFloatToInt(t *testing.T) {
	str := "12.13"

	val, err := grinex.StrFloatToInt(str, 2)
	if err != nil {
		t.Error(err)
	}

	if val != 1213 {
		t.Errorf("wrong conversion, want 1213, got:%d", val)
	}
}
