package filePool

import (
	"testing"
)

var testLen = []struct {
	len    int
	expect int
}{
	{0, 0},
	{1, 1},
	{10, 10},
	{15, 15},
	{100, 100},
}

const MAX = 100

func TestFilePool(t *testing.T) {

	for _, i := range testLen {
		content := ""
		var target []string
		fs := FileSetNew()
		for j := 0; j < i.len; j++ {
			content = content + "a"
			fs.Add(content)
			target = append(target, content)
		}
		if fs.Length() == i.expect {
			t.Logf("test: len %d Add method Pass!", i)
		} else {
			t.Errorf("test: len %d Add method NO PASS !!", i)
		}
		n := fs.Length()
		//for i := 0; i < n; i++ {
		pass := true
		for i := 0; i < n; i++ {
			s := fs.Get()
			if s != target[i] {
				pass = false
			}
		}
		if pass {
			t.Logf("test :Get() %d  succuss!! Pass!", i)
		} else {
			t.Errorf("test:Get() %d Failed!! NO PASS! ", i)
		}
	}
}
