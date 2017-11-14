package hld

import (
	"testing"
)

func BenchmarkChannelSync(b *testing.B) {
	ch := make(chan int)
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		close(ch)
	}()
	for _ = range ch {
	}
}

func TestAverages(t *testing.T) {
	if !FileExists("./test/utf8.txt") {
		t.Error("test text file is missing")
	}

	var data string
	var err error
	data, err = HSFileToString("./test/utf8.txt")
	if err != nil {
		t.Error(err.Error())
	}
	if data != HSTestString {
		t.Error("utf8 failed")
	}

	data, err = HSFileToString("./test/utf16.txt")
	if err != nil {
		t.Error(err.Error())
	}
	if data != HSTestString {
		t.Error("utf16 failed")
	}

	data, err = HSFileToString("./test/utf16b.txt")
	if err != nil {
		t.Error(err.Error())
	}
	if data != HSTestString {
		t.Error("utf16 BE failed")
	}

}
