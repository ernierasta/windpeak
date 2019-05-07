package hashbench

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/OneOfOne/xxhash"
)

func getMD5(path string, info os.FileInfo) (string, error) {
	fmd5 := ""
	if !info.IsDir() && filepath.Ext(path) != ".bsa" { // skip md5 for bsa
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()

		h := md5.New()

		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}
		fmd5 = fmt.Sprintf("%x", h.Sum(nil))
	}
	return fmd5, nil
}

// we need to use 32-bit hashing for 32-bit platforms!
func getXXhash(path string, info os.FileInfo) (string, error) {
	fxx := ""
	if !info.IsDir() && filepath.Ext(path) != ".bsa" { // skip hash for bsa
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()

		h := xxhash.New64()

		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}
		fxx = fmt.Sprintf("%v", h.Sum64())
	}
	return fxx, nil
}
