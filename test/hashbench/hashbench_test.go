package hashbench

import (
	"io/ioutil"
	"os"
	"testing"
)

var bigf = "delme.test"

func setup() {
	bigBuff := make([]byte, 750000000)
	ioutil.WriteFile(bigf, bigBuff, 0666)
}

func cleanup() {
	os.Remove(bigf)
}

func BenchmarkMD5(b *testing.B) {
	setup()
	info, _ := os.Stat(bigf)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		getMD5(bigf, info)
	}
}

func BenchmarkXXhash(b *testing.B) {
	info, _ := os.Stat(bigf)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		getXXhash(bigf, info)
	}
}
