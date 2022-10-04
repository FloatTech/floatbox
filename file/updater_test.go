package file

import (
	"os"
	"testing"
)

func TestLazy(t *testing.T) {
	err := os.MkdirAll("data/Control", 0755)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetLazyData("data/Control/kanban.png", true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCustLazy(t *testing.T) {
	err := os.MkdirAll("data/Tarot", 0755)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetCustomLazyData("https://gitcode.net/shudorcl/zbp-tarot/-/raw/master/", "data/Tarot/tarots.json")
	if err != nil {
		t.Fatal(err)
	}
}
