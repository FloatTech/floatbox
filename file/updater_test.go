package file

import "testing"

func TestLazy(t *testing.T) {
	_, err := GetLazyData("data/Control/kanban.png", true)
	if err != nil {
		t.Fatal(err)
	}
}
