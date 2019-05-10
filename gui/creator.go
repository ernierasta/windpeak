package gui

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/pkg/browser"
	"github.com/sqweek/dialog"
)

const (
	NO_ID       = -11111
	MP_ID       = 50000
	SMP_BASE_ID = 10000
)

type Gui struct {
	settings         *Settings
	w                *ui.Window
	tabs             *DynamicTab
	modsTabID        int
	stepOptionsTabID int
	stepModsButtons  []*ui.Button
	stepListGrp      *ui.Group
	stepListVbox     *DynamicBox

	mpTab map[string]Input
	smTab map[string]interface{}

	cbLoadAllFunc  func(dirpath string) (*Modpack, error)
	cbSaveAllFunc  func(mp *Modpack, dirpath string) error
	cbCmprFileFunc func(mp *Modpack, file string) error
	cbSaveLOFunc   func(s *Step) error
	cbNewModFunc   func() error

	saved bool
}

func NewCreatorGUI(s *Settings) *Gui {
	g := &Gui{settings: s}
	// initialize Callbacks, so code will not panic
	g.SetCBLoadAll(func(p string) (mp *Modpack, err error) { return })
	g.SetCBSaveAll(func(mp *Modpack, dirpath string) (err error) { return })
	g.SetCBCompressAllToFile(func(mp *Modpack, file string) (err error) { return })
	g.SetCBSaveStepLO(func(s *Step) (err error) { return })
	g.SetCBCreateNewMod(func() (err error) { return })
	return g
}

func (g *Gui) Show() {
	ui.Main(g.show)
}

func (g *Gui) SetCBLoadAll(f func(path string) (*Modpack, error)) {
	g.cbLoadAllFunc = f
}

func (g *Gui) SetCBSaveAll(f func(mp *Modpack, path string) error) {
	g.cbSaveAllFunc = f
}

func (g *Gui) SetCBCompressAllToFile(f func(mp *Modpack, file string) error) {
	g.cbCmprFileFunc = f
}

func (g *Gui) SetCBCreateNewMod(f func() error) {
	g.cbNewModFunc = f
}
func (g *Gui) SetCBSaveStepLO(f func(s *Step) error) {
	g.cbSaveLOFunc = f
}

func (g *Gui) show() {
	g.w = ui.NewWindow("Windpeak Modpack Creator Mode", 640, 480, true)
	g.w.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		g.w.Destroy()
		return true
	})

	g.tabs = NewDynamicTab()
	g.w.SetChild(g.tabs)
	g.w.SetMargined(true)

	g.tabs.Append("Mod manager", g.modmanagerTab(), MP_ID+1)
	g.tabs.SetMargined(MP_ID+1, true)
	g.tabs.Append("Modpack", g.modpackTab(g.settings.Modpack), MP_ID)
	g.tabs.SetMargined(MP_ID, true)

	g.w.Show()
}

func (g *Gui) modmanagerTab() ui.Control {
	grp := ui.NewGroup("Select mod manager you want to use:")
	grp.SetMargined(true)
	grid := ui.NewGrid()
	grid.SetPadded(true)
	grp.SetChild(grid)
	radio := ui.NewRadioButtons()
	radio.Append("Mod Organizer 2")
	radio.Append("Wrye Bash")
	grid.Append(ui.NewHorizontalSeparator(), 0, 0, 1, 1, true, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(ui.NewHorizontalSeparator(), 1, 0, 1, 1, false, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(ui.NewHorizontalSeparator(), 2, 0, 1, 1, true, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(ui.NewHorizontalSeparator(), 0, 1, 1, 1, true, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(radio, 1, 1, 1, 1, false, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(ui.NewHorizontalSeparator(), 2, 1, 1, 1, true, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(ui.NewHorizontalSeparator(), 0, 2, 1, 1, true, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(ui.NewButton("Choose"), 1, 2, 1, 1, false, ui.AlignCenter, false, ui.AlignCenter)
	grid.Append(ui.NewHorizontalSeparator(), 2, 2, 1, 1, true, ui.AlignCenter, false, ui.AlignCenter)

	return grp
}

// TODO: this is useless until libui will have support for images
func (g *Gui) loadLogo() *ui.Image {
	file := "assets/Windpeak_Inn_Shop_Sign.gif"
	logo, err := os.Open(file)
	if err != nil {
		ui.MsgBoxError(g.w, "Error", "Can not open logo image from: "+file+", err: "+err.Error())
		return &ui.Image{}
	}
	defer logo.Close()

	img, _, err := image.Decode(logo)
	if err != nil {
		ui.MsgBoxError(g.w, "Error", "Can not decode logo image, err:"+err.Error())
		return &ui.Image{}
	}

	rect := img.Bounds()
	logoRGBA := image.NewRGBA(rect)
	draw.Draw(logoRGBA, rect, img, rect.Min, draw.Src)

	i := ui.NewImage(437, 512)
	i.Append(logoRGBA)
	return i
}

func (g *Gui) modpackTab(mp *Modpack) ui.Control {

	g.mpTab = map[string]Input{"DirPath": ui.NewLabel(""), "Name": ui.NewEntry(), "Version": ui.NewEntry(),
		"Author": ui.NewEntry(), "Homepage": ui.NewEntry(), "Description": ui.NewMultilineEntry()}
	g.modpackTabUpdate()

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	grp := ui.NewGroup("Modpack metadata")
	f := ui.NewForm()
	f.SetPadded(true)
	f.Append("Dir", g.mpTab["DirPath"], false)
	f.Append("Name", g.mpTab["Name"], false)
	f.Append("Version", g.mpTab["Version"], false)
	f.Append("Author", g.mpTab["Author"], false)
	f.Append("Homepage", g.mpTab["Homepage"], false)
	f.Append("Description", g.mpTab["Description"], false)
	grp.SetChild(f)
	grp.SetMargined(true)
	vbox.Append(grp, false)

	vbox.Append(g.modpackControls(), false)
	vbox.Append(ui.NewHorizontalSeparator(), false)
	vbox.Append(g.stepList(), false)

	return vbox
}

func (g *Gui) modpackTabUpdate() {
	g.modpackOptionsUpdate()
}

func (g *Gui) modpackOptionsUpdate() {
	g.mpTab["DirPath"].SetText(g.settings.Modpack.DirPath)
	g.mpTab["Name"].SetText(g.settings.Modpack.Name)
	g.mpTab["Version"].SetText(g.settings.Modpack.Version)
	g.mpTab["Author"].SetText(g.settings.Modpack.Author)
	g.mpTab["Homepage"].SetText(g.settings.Modpack.Homepage)
	g.mpTab["Description"].SetText(g.settings.Modpack.Description)
}

// modsTab returns button callback which creates new tab containing step mods
func (g *Gui) modsTab(s *Step) func(b *ui.Button) {
	return func(b *ui.Button) {
		//remove previously opened step tab & enable all buttons
		if g.modsTabID != 0 {
			g.tabs.Delete(g.modsTabID)
			for i := range g.stepModsButtons {
				g.stepModsButtons[i].Enable()
			}
		}
		vbox := NewDynamicVbox()
		for i := range s.Mods {
			hbox := ui.NewHorizontalBox()
			hbox.SetPadded(true)
			b := ui.NewButton("Homepage")
			b.OnClicked(func(b *ui.Button) {
				browser.OpenURL(s.Mods[i].HomepageURL)
			})
			hbox.Append(ui.NewLabel(fmt.Sprintf("%d. %s", i, s.Mods[i].Name)), true)
			hbox.Append(b, false)
			vbox.Append(hbox, s.Mods[i].ID)
		}
		vbox.Append(g.modsControls(s), NO_ID)
		g.modsTabID = s.ID
		g.tabs.Append(fmt.Sprintf("Step: %d. %s", s.ID, s.Name), vbox, s.ID)
		g.tabs.SetMargined(s.ID, true)
		b.Disable()
	}
}

func (g *Gui) stepOptionsTab(s *Step) {

	if g.stepOptionsTabID != 0 {
		g.tabs.Delete(g.stepOptionsTabID)
	}
	g.stepOptionsTabID = s.ID + SMP_BASE_ID

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	f := ui.NewForm()
	f.SetPadded(true)
	grp := ui.NewGroup("Step settings")
	grp.SetMargined(true)
	vbox.Append(f, false)
	grp.SetChild(vbox)

	g.smTab = map[string]interface{}{
		"ID":            ui.NewEntry(),
		"Name":          ui.NewEntry(),
		"Description":   ui.NewMultilineEntry(),
		"AlternativeTo": ui.NewSpinbox(0, 1000),
		"Optional":      ui.NewCheckbox("Optional"),
		"Stop":          ui.NewCheckbox("Stop processing"),
		"StopMsg":       ui.NewMultilineEntry(),
		"StopCommands":  ui.NewMultilineEntry(),
	}

	f.Append("ID", g.smTab["ID"].(*ui.Entry), false)
	f.Append("Name", g.smTab["Name"].(*ui.Entry), false)
	f.Append("Description", g.smTab["Description"].(*ui.MultilineEntry), true)
	f.Append("Alternative To step with ID", g.smTab["AlternativeTo"].(*ui.Spinbox), false)
	f.Append("Optional step", g.smTab["Optional"].(*ui.Checkbox), false)
	f.Append("Stop processing after this step is finished", g.smTab["Stop"].(*ui.Checkbox), false)
	f.Append("Message shown to user on step stop", g.smTab["StopMsg"].(*ui.MultilineEntry), true)
	f.Append("StopCommands", g.smTab["StopCommands"].(*ui.MultilineEntry), true)

	g.stepOptionsTabUpdate(s)
	g.smTab["ID"].(*ui.Entry).Disable()
	save := ui.NewButton("Save")
	save.OnClicked(func(b *ui.Button) {
		s := g.stepOptionsGUIRead()
		// g.stepSaveCB(s)
		err := g.stepSaveToModpack(s)
		if err != nil {
			ui.MsgBoxError(g.w, "error", err.Error())
			return
		}
		g.stepListUpdate()
		g.tabs.Delete(g.stepOptionsTabID)
		g.stepOptionsTabID = 0
	})
	cancel := ui.NewButton("Cancel")
	cancel.OnClicked(func(b *ui.Button) {
		g.tabs.Delete(g.stepOptionsTabID)
		g.stepOptionsTabID = 0
	})
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	hbox.Append(save, false)
	hbox.Append(cancel, false)
	vbox.Append(hbox, false)
	g.tabs.Append(fmt.Sprintf("Options: %d. %s", s.ID, s.Name), grp, s.ID+SMP_BASE_ID)
	g.tabs.SetMargined(s.ID+SMP_BASE_ID, true)
}

func (g *Gui) stepOptionsTabUpdate(s *Step) {
	g.smTab["ID"].(*ui.Entry).SetText(fmt.Sprintf("%d", s.ID))
	g.smTab["Name"].(*ui.Entry).SetText(s.Name)
	g.smTab["Description"].(*ui.MultilineEntry).SetText(s.Description)
	g.smTab["AlternativeTo"].(*ui.Spinbox).SetValue(s.AlternativeTo)
	g.smTab["Optional"].(*ui.Checkbox).SetChecked(s.Optional)
	g.smTab["Stop"].(*ui.Checkbox).SetChecked(s.Stop)
	g.smTab["StopMsg"].(*ui.MultilineEntry).SetText(s.StopMsg)
	g.smTab["StopCommands"].(*ui.MultilineEntry).SetText(strings.Join(s.StopCommands, "\n"))
}

func (g *Gui) stepOptionsGUIRead() *Step {
	s := &Step{}
	id := g.smTab["ID"].(*ui.Entry).Text()
	s.ID, _ = strconv.Atoi(id)
	s.Name = g.smTab["Name"].(*ui.Entry).Text()
	s.Description = g.smTab["Description"].(*ui.MultilineEntry).Text()
	s.AlternativeTo = g.smTab["AlternativeTo"].(*ui.Spinbox).Value()
	s.Optional = g.smTab["Optional"].(*ui.Checkbox).Checked()
	s.Stop = g.smTab["Stop"].(*ui.Checkbox).Checked()
	s.StopMsg = g.smTab["StopMsg"].(*ui.MultilineEntry).Text()
	s.StopCommands = strings.Split(g.smTab["StopCommands"].(*ui.MultilineEntry).Text(), "\n")

	return s
}

func (g *Gui) stepSaveToModpack(s *Step) error {

	max := 0
	for i := range g.settings.Modpack.Steps {
		if g.settings.Modpack.Steps[i].ID > max {
			max = g.settings.Modpack.Steps[i].ID
		}
		if g.settings.Modpack.Steps[i].ID == s.ID {
			g.settings.Modpack.Steps[i] = s
			return nil
		}
	}
	if s.ID == 0 {
		s.ID = max + 1
		g.settings.Modpack.Steps = append(g.settings.Modpack.Steps, s)
		return nil
	}
	return fmt.Errorf("wrong Step.ID value: %d, not found and not new", s.ID)
}

func (g *Gui) stepList() ui.Control {
	g.stepListGrp = ui.NewGroup("Steps")
	g.stepListGrp.SetMargined(true)
	g.stepListCreate()
	g.stepListGrp.SetChild(g.stepListVbox)
	return g.stepListGrp
}

func (g *Gui) stepListCreate() {
	g.stepListVbox = NewDynamicVbox()
	for i := range g.settings.Modpack.Steps {
		g.stepListVbox.Append(g.stepListLine(g.settings.Modpack.Steps[i]), g.settings.Modpack.Steps[i].ID)
	}
	g.stepListVbox.Append(g.stepListControls(), NO_ID)
}

func (g *Gui) stepListUpdate() {
	g.stepListVbox.DeleteAll()
	g.stepListCreate()
	g.stepListGrp.SetChild(g.stepListVbox)

}

func (g *Gui) stepListLine(s *Step) ui.Control {
	stepModsB := ui.NewButton("⚔ Mods ⚔")
	g.stepModsButtons = append(g.stepModsButtons, stepModsB)
	stepModsB.OnClicked(g.modsTab(s))
	hbox := ui.NewHorizontalBox()
	hbox.Append(ui.NewLabel(s.Name), true)
	hbox.Append(stepModsB, false)
	//hbox.Append(ui.NewButton("↑"), false)
	//hbox.Append(ui.NewButton("↓"), false)
	return hbox
}

func (g *Gui) stepListControls() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	n := ui.NewButton("⛊ New")
	n.OnClicked(func(b *ui.Button) {
		g.stepOptionsTab(&Step{})
	})
	c := ui.NewCombobox()
	c.Hide()
	names := []string{}
	for i := range g.settings.Modpack.Steps {
		names = append(names, g.settings.Modpack.Steps[i].Name)
	}
	e := NewComboButton("⛊ Edit", names, c)
	e.OnClicked(func(b *ui.Button) {
		g.stepOptionsTab(g.settings.Modpack.Steps[e.Selected()])
	})

	hbox.Append(n, false)
	hbox.Append(e, false)
	hbox.Append(c, true)
	return hbox
}

func (g *Gui) modsControls(s *Step) ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	n := ui.NewButton("প New")
	n.OnClicked(func(b *ui.Button) {
		err := g.cbNewModFunc()
		if err != nil {
			ui.MsgBoxError(g.w, "Error", "Error occured while creating mod definitions: "+err.Error())
		}
	})
	e := ui.NewButton("Save plugins & LO")
	e.OnClicked(func(b *ui.Button) {
		err := g.cbSaveLOFunc(s)
		if err != nil {
			ui.MsgBoxError(g.w, "Error", "Error occured writing load orded or plugins list: "+err.Error())
		}
	})
	hbox.Append(n, false)
	hbox.Append(e, false)
	return hbox
}

func (g *Gui) modpackControls() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	g.modpackButtonNew(hbox)
	g.modpackButtonOpen(hbox)
	g.modpackButtonSave(hbox)
	g.modpackButtonSaveAs(hbox)
	return hbox
}

func (g *Gui) modpackButtonNew(hbox *ui.Box) {
	b := ui.NewButton("❋ New")
	b.OnClicked(func(button *ui.Button) {
		ok := dialog.Message("%s", "Do you want to create new modpack? Steps are are saved automaticaly.").
			Title("New modpack").YesNo()
		if ok {
			g.settings.Modpack = &Modpack{}
			g.modpackTabUpdate()
		}
	})
	hbox.Append(b, false)
}

func (g *Gui) modpackButtonOpen(hbox *ui.Box) {
	b := ui.NewButton("Open")
	b.OnClicked(func(button *ui.Button) {
		var ok = true
		if !g.saved {
			ok = dialog.Message("%s\n\n%s", "Do you want to open new modpack?", "Current will be discarded!").
				Title("Open modpack").YesNo()
		}
		if ok {
			d := dialog.Directory().Title("Select modpack directory")
			exe, err := os.Executable()
			if err != nil {
				ui.MsgBoxError(g.w, "Error", "Can't get executable filename, err: "+err.Error())
				return
			}

			d.StartDir = filepath.Dir(exe)
			dir, err := d.Browse()
			if err == nil {
				g.cbLoadAllFunc(dir)
				g.modpackTabUpdate()
			}
		}
	})

	hbox.Append(b, false)
}

func (g *Gui) modpackButtonSave(hbox *ui.Box) {
	b := ui.NewButton("Save")
	b.OnClicked(func(button *ui.Button) {
		err := g.cbSaveAllFunc(g.settings.Modpack, g.settings.Modpack.DirPath)
		if err != nil {
			ui.MsgBoxError(g.w, "Error", "Can't save modpack, err: "+err.Error())
			return
		}
		g.modpackTabUpdate()
	})
	hbox.Append(b, false)
}

func (g *Gui) modpackButtonSaveAs(hbox *ui.Box) {
	b := ui.NewButton("Save as ❋")
	b.OnClicked(func(button *ui.Button) {
		d := dialog.Directory().Title("Save as modpack directory")
		exe, err := os.Executable()
		if err != nil {
			ui.MsgBoxError(g.w, "Error", "Can't get executable filename, err: "+err.Error())
			return
		}

		d.StartDir = filepath.Dir(exe)
		dir, err := d.Browse()
		if err == nil {
			err := g.cbSaveAllFunc(g.settings.Modpack, dir)
			if err != nil {
				ui.MsgBoxError(g.w, "Error", "Can't save mod into dir: "+dir+", err: "+err.Error())
				return
			}
			g.modpackTabUpdate()
		}
	})

	hbox.Append(b, false)
}
