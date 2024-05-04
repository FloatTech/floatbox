package file

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLazy(t *testing.T) {
	err := os.MkdirAll("data/Control", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll("data/control", 0755)
	if err != nil {
		t.Fatal(err)
	}
	data, err := GetLazyData("data/Control/kanban.png", "data/control/stor.spb", true)
	if err != nil {
		t.Fatal(err)
	}
	b := md5.Sum(data)
	assert.Equal(t, "806dedf39de375531f155285cf6ddb61", hex.EncodeToString(b[:]))
}

func TestCustLazy(t *testing.T) {
	err := os.MkdirAll("data/Tarot", 0755)
	if err != nil {
		t.Fatal(err)
	}
	data, err := GetCustomLazyData("https://gitcode.net/shudorcl/zbp-tarot/-/raw/master/", "data/Tarot/tarots.json")
	if err != nil {
		t.Fatal(err)
	}
	b := md5.Sum(data)
	assert.Equal(t, "e0aa021e21dddbd6d8cecec71e9cf564", hex.EncodeToString(b[:]))
}
