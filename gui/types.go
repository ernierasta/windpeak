package gui

import "github.com/andlabs/ui"

type Settings struct {
	Modpack              *Modpack
	DownloadDir, GameDir string
}

// Input allows to abstract text input Entries
type Input interface {
	ui.Control
	SetText(text string)
}

type Modpack struct {
	DirPath               string
	Name, Version, Author string
	Homepage, Description string
	InstallFolder         string
	Steps                 []*Step
}

// Step represents one Modpack step
type Step struct {
	ID             int
	Name           string
	Description    string
	AlternativeTo  int
	Optional, Stop bool
	StopMsg        string
	StopCommands   []string
	Mods           []*Mod
	// TODO: here or elswere:
	// Actions []Action
	// copy files from main step dir
}

type Mod struct {
	ID                       int
	Name, Version, Author    string
	HomepageURL, DownloadURL string
}
