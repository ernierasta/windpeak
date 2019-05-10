package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/icedream/binarydist"
	"github.com/pquerna/snapdiff"
	//"github.com/monmohan/xferspdy"

	"github.com/pietroglyph/xferspdy"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// let's test diff and patch functions

func main() {

	//basefile := "sr_Oblivion_Stutter_Remover"
	basefile := "Oblivion"
	f1, err := ioutil.ReadFile(basefile + ".org.ini")
	if err != nil {
		log.Fatal(err)
	}
	f2, err := ioutil.ReadFile(basefile + ".mod.ini")
	if err != nil {
		log.Fatal(err)
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(f1), string(f2), true)
	patch := dmp.PatchMake(diffs)
	fmt.Println("patch: ", patch)
	patchStr := dmp.PatchToText(patch)
	fmt.Println("patch string: ", patchStr)
	ioutil.WriteFile("my.patch", []byte(patchStr), 0644)
	patchFromFile, err := ioutil.ReadFile("my.patch")
	if err != nil {
		log.Fatal(err)
	}
	patch2, err := dmp.PatchFromText(string(patchFromFile))
	if err != nil {
		log.Fatal(err)
	}
	t22, applied := dmp.PatchApply(patch2, string(f1))
	_ = t22
	fmt.Println("applied succ: ", applied)
	//fmt.Println(t22)

	fmt.Println("apply patches to different Oblivion.ini")
	f3, _ := ioutil.ReadFile(basefile + ".different.ini")
	t33, applied33 := dmp.PatchApply(patch2, string(f3))
	fmt.Println(t33)
	fmt.Println(applied33)

	// test bin patching
	fst := "/home/ernie/Temp/delme/Dragonborn.original.esm"
	sec := "/home/ernie/Temp/delme/Dragonborn.cleaned.esm"
	des := "."
	_, _ = sec, des
	bdStart := time.Now()
	fmt.Println(binMakePatchFileBinaryDist(fst, sec, des))
	fmt.Printf("dinarydist time: %f\n", time.Since(bdStart).Seconds())
	patchStart := time.Now()
	fmt.Println(binDoPatch(fst, "Dragonborn.original.esm.binpatch", des))
	fmt.Printf("patch time: %f\n", time.Since(patchStart).Seconds())

}

// This will be used!
func binMakePatchFileBinaryDist(fst, sec, destdir string) error {
	_, fname := filepath.Split(fst)
	fstFile, err := os.Open(fst)
	if err != nil {
		return err
	}
	defer fstFile.Close()
	secFile, err := os.Open(sec)
	if err != nil {
		return err
	}
	defer secFile.Close()
	patchFile, err := os.OpenFile(filepath.Join(destdir, fname+".binpatch"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer patchFile.Close()

	return binarydist.Diff(fstFile, secFile, patchFile)
}

func binDoPatch(org, patch, destdir string) error {
	_, fname := filepath.Split(org)
	orgFile, err := os.Open(org)
	if err != nil {
		return err
	}
	defer orgFile.Close()
	patchFile, err := os.Open(patch)
	if err != nil {
		return err
	}
	defer patchFile.Close()
	outputFile, err := os.OpenFile(filepath.Join(destdir, fname+".new"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer patchFile.Close()

	return binarydist.Patch(orgFile, outputFile, patchFile)

}

func binMakePatchFileSnappy(fst, sec, destdir string) error {
	_, fname := filepath.Split(fst)
	fstFile, err := os.Open(fst)
	if err != nil {
		return err
	}
	defer fstFile.Close()
	secFile, err := os.Open(sec)
	if err != nil {
		return err
	}
	defer secFile.Close()
	patchFile, err := os.OpenFile(filepath.Join(destdir, fname+".snappybinpatch"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer patchFile.Close()

	return snapdiff.Diff(fstFile, secFile, patchFile)

}

func binMakeFingerprintFile(file, destdir string) error {
	fingerprint, err := xferspdy.NewFingerprint(file, 2*1024)
	if err != nil {
		return err
	}

	dir, filename := filepath.Split(file)
	fpFilepath := filepath.Join(dir, filename+".fingerprint")
	fpfile, err := os.OpenFile(fpFilepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(fpfile)
	enc.Encode(*fingerprint)
	fpfile.Close()
	return nil
}

func binMakePatchFile(fst, sec, destdir string) error {
	fstFingerprint, err := xferspdy.NewFingerprint(fst, 2*1024)
	if err != nil {
		return err
	}
	fmt.Println("fp done!")
	diff, err := xferspdy.NewDiff(sec, *fstFingerprint)
	if err != nil {
		return err
	}
	fmt.Println("diff done!")
	_, fname := filepath.Split(fst)
	patchFile, err := os.OpenFile(filepath.Join(destdir, fname+".binpatch"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer patchFile.Close()

	enc := gob.NewEncoder(patchFile)
	err = enc.Encode(diff)
	if err != nil {
		return err
	}
	return nil
}
