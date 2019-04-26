package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sergi/go-diff/diffmatchpatch"

	home "github.com/mitchellh/go-homedir"
)

func main() {
	beforeData, err := getDirStructure(`C:\Games\Oblivion\Data`)
	if err != nil {
		fmt.Println(err)
	}
	beforeModsDir, err := getDirStructure(`C:\Games\Oblivion mods`)
	if err != nil {
		fmt.Println(err)
	}
	beforeHomeDir, err := getDirStructure(`C:\Users\ernie\\Documents\My Games\Oblivion\`)
	if err != nil {
		fmt.Println(err)
	}
	beforeini, err := getFile()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Now install mod")
	reader := bufio.NewReader(os.Stdin)
	_, _, err = reader.ReadRune()

	afterData, err := getDirStructure(`C:\Games\Oblivion\Data`)
	if err != nil {
		fmt.Println(err)
	}
	afterModsDir, err := getDirStructure(`C:\Games\Oblivion mods`)
	if err != nil {
		fmt.Println(err)
	}
	afterHomeDir, err := getDirStructure(`C:\Users\ernie\\Documents\My Games\Oblivion\`)
	if err != nil {
		fmt.Println(err)
	}
	afterini, err := getFile()
	if err != nil {
		fmt.Println(err)
	}

	removedDD := difference(beforeData, afterData)
	addedDD := difference(afterData, beforeData)
	removedMD := difference(beforeModsDir, afterModsDir)
	addedMD := difference(afterModsDir, beforeModsDir)
	changesHD := difference(afterHomeDir, beforeHomeDir)
	p(addedDD, "DATA Dir added:")
	p(removedDD, "DATA Dir removed:")
	p(addedMD, "MOD DIR added:")
	p(removedMD, "MOD DIR removed:")
	p(changesHD, "Home changes:")

	//fmt.Println("before:", beforeModsDir)
	//fmt.Println("after:", afterModsDir)

	dmp := diffmatchpatch.New()
	diff := dmp.DiffMain(beforeini, afterini, false)
	fmt.Println("Changes in ini:")
	fmt.Println(DiffPrettyText(diff))

}

func p(files map[string]string, title string) {
	fmt.Println(title)
	for file, id := range files {
		fmt.Printf("%q: %v\n", file, id)
	}
}

func getDirStructure(dir string) (map[string]string, error) {
	files := map[string]string{}
	return files, filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//files[path] = getMD5(path, info)
		files[path] = makeID(info)
		return nil
	})
}

func makeID(info os.FileInfo) string {
	return strconv.FormatInt(info.ModTime().Unix()+info.Size(), 10)
}

func getMD5(path string, info os.FileInfo) (string, error) {
	fmd5 := ""
	if !info.IsDir() && filepath.Ext(path) != "bsa" { // skip md5 for bsa
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

func getFile() (string, error) {
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

func difference(a, b map[string]string) map[string]string {

	ret := map[string]string{}
	for ka, va := range a {
		if vb, ok := b[ka]; !ok || va != vb { // if key is missing or md5 different
			ret[ka] = vb
		}
	}
	return ret
}

// DiffPrettyText converts a []Diff into a colored text report.
func DiffPrettyText(diffs []diffmatchpatch.Diff) string {
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
