package main

import (
	"github.com/andlabs/ui"
)

// this is draft, how Gui gui and methods could look
// for now ignored ;-)

const (
	NO_ID       = -11111
	MP_ID       = 50000
	SMP_BASE_ID = 10000
)

type Gui struct {
	settings *Settings
	w        *ui.Window
	tabs     *gui.DynamicTab
}

func NewCreatorGUI(s *Settings) *Gui {
	g := &Gui{settings: s}
	// initialize Callbacks, so code will not panic
	g.SetModpackLoadCallback(func(p string) (mp *gui.Modpack, err error) { return })
	//g.SetModpackSaveCallback(func(mp *gui.Modpack, dirpath string) (err error) { return })
	//g.SetModpackExportCallback(func(mp *gui.Modpack, file string) (err error) { return })
	//g.SetStepSaveLoCallback(func(s *gui.Step) (err error) { return })
	//g.SetModCreateCallback(func() (err error) { return })
	return g
}

func (g *Gui) Show() {}
func (g *Gui) SetCBLoadAll(f func(path string) (*Modpack, error))        {}
func (g *Gui) SetCBSaveAll(f func(mp *Modpack, path string) error)       {}
func (g *Gui) SetCBSaveAllToArchive(f func(mp *Modpack, file string) error) {}
func (g *Gui) SetCBSaveStepLO(f func(s *gui.Step) (err error) {}
func (g *Gui) SetCBCreateMod(f func() (err error)) {}

func (g *Gui) append(t *ui.Control) {}

type ModManagerTab struct{}

type ModpackTab struct{}

func (mpt *ModpackTab) update(mp *Modpack) {}

type StepOptionsTab struct{}

type StepModsTab struct{}
