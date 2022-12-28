package file

import (
	"hash/crc64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDl(t *testing.T) {
	_ = os.RemoveAll("啊")
	err := os.MkdirAll("啊", 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("啊")
	err = DownloadTo("https://gitcode.net/u011570312/zbpdata/-/raw/main/Control/kanban.png", "啊/看板.png")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("啊/看板.png")
	if err != nil {
		t.Fatal(err)
	}
	h := crc64.New(crc64.MakeTable(crc64.ECMA))
	_, err = h.Write(data)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint64(0x814b82a8efcbcef2), h.Sum64())
}
