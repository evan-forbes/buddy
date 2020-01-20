package abigen

import "testing"

func TestFindFile(t *testing.T) {
	path := "/etc"
	ext := "group"
	_, has := findFile(path, ext)
	if !has {
		t.Error("could not find file name 'group' in /etc/")
	}
}
