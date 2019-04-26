package main

import (
	"fmt"
)

type ModPack struct {
	Name          string
	Version       string
	Author        string
	SourceUrl     string
	Description   string
	TargetGame    string
	InstallFolder string
}

func (mp *ModPack) Read() error {
	return nil
}

func (mp *ModPack) InstallMods() error {
	return nil
}

type Mod struct {
	ID             int64
	Name           string
	Version        string
	FileName       string
	MD5            string
	Url            string
	Repository     string
	Description    string
	FileID         int64
	InstalledFiles []string
}

func (m *Mod) Install() error {
	return nil
}

func main() {
	m := Mod{}
	m.ReadMeta("OBSE Test Plugin-33574.rar.meta", "OBSE Tester.ini")

	fmt.Printf("%+v\n", m)
}
