package mod

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/OneOfOne/xxhash"
	home "github.com/mitchellh/go-homedir"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Rule struct {
	Source      string
	Destination string
}

type Mod struct {
	ID                       int
	Name, Version, Author    string
	HomepageURL, DownloadURL string
	InstallRules             []*Rule
}

type WryeBashMod struct {
	dowloadDir string
	gameDir    string
}

func NewWreyBashMod(downloadDir, gameDir string) *WryeBashMod {
	return &WryeBashMod{
		dowloadDir: downloadDir,
		gameDir:    gameDir,
	}
}

// ReadMeta reads meta files:
// - meta.ini
func (wbm *WryeBashMod) ReadMeta() *Mod {
	return &Mod{}
}

// WriteMeta writes all meta files.
// - meta.ini - keeps info about mod and install rules
func (wbm *WryeBashMod) WriteMeta(m *Mod) error {
	return nil
}

// CreateMeta retrieves mod data from two sources:
// 1. Meta files if availabile
// 2. Snapshoting game directory and comparing changes
func (wbm *WryeBashMod) CreateMeta(fast, xxhash bool) (*Mod, error) {
	gameSnap := NewSnapshot(wbm.gameDir)
	gameSnap.Make(fast, xxhash)
	return &Mod{}, nil
}

// TODO: move snapshot to separate file, MO2 will use it also (archive -> moddir diff)
type Snapshot struct {
	dir     string
	first   map[string]string
	archive map[string]string
	final   map[string]string
}

func NewSnapshot(dir string) *Snapshot {
	return &Snapshot{dir: dir}
}

func (st *Snapshot) Make(fast bool, xxhash bool) (map[string]string, error) {
	files := map[string]string{}
	return files, filepath.Walk(st.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fast {
			files[path] = st.makeID(info)
		} else {
			var err error
			if xxhash {
				files[path], err = st.getXXhash(path, info)
			} else {
				files[path], err = st.getMD5(path, info)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})

}
func (st *Snapshot) Archive(file string) (m map[string]string, err error) { return }
func (st *Snapshot) Diff()                                                {}

func (st *Snapshot) makeID(info os.FileInfo) string {
	if info.IsDir() { // skip dirs
		return ""
	}
	return strconv.FormatInt(info.ModTime().Unix()+info.Size(), 10)
}

func (st *Snapshot) getMD5(path string, info os.FileInfo) (string, error) {
	fmd5 := ""
	if !info.IsDir() {
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
func (st *Snapshot) getXXhash(path string, info os.FileInfo) (string, error) {
	fxx := ""
	if !info.IsDir() {
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

func (wbm *WryeBashMod) getFile() (string, error) {
	h, err := home.Dir()
	if err != nil {
		return "", err
	}
	f, err := ioutil.ReadFile(h + `\Documents\My Games\Oblivion\Oblivion.ini`)
	if err != nil {
		return "", nil
	}
	return string(f), nil
}

func (wbm *WryeBashMod) difference(a, b map[string]string) map[string]string {

	ret := map[string]string{}
	for ka, va := range a {
		// return if key is missing or key is there, but value is different for files other then esp or esm
		// we are avoiding returning files changed becouse of load order reorganization
		if vb, ok := b[ka]; !ok || (va != vb && filepath.Ext(ka) != ".esp" && filepath.Ext(ka) != ".esm") { // if key is missing or value different
			ret[ka] = vb
		}
	}
	return ret
}

// DiffPrettyText converts a []Diff into a colored text report.
func (wbm *WryeBashMod) DiffPrettyText(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString("\x1b[32m")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("\x1b[0m")
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString("\x1b[31m")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("\x1b[0m")
		case diffmatchpatch.DiffEqual:
			// do not show equals
		}
	}

	return buff.String()
}

func (wbm *WryeBashMod) getLoadOrder(dir string) ([]os.FileInfo, error) {
	esps := []os.FileInfo{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != dir {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".esp" || filepath.Ext(path) == ".esm" {
			esps = append(esps, info)
		}
		return nil
	})

	sort.Slice(esps, func(i, j int) bool {
		return esps[i].ModTime().Unix() < esps[j].ModTime().Unix()
	})
	return esps, err
}
