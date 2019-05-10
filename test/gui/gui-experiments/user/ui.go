package main

import (
	"fmt"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/pkg/browser"
)

var (
	w                *ui.Window
	items, installed []*Mod
	mp               *Modpack
)

type Mod struct {
	ID                         int
	Name, Author, Version, URL string
}

type Modpack struct {
	Name, Version, Author  string
	SourceUrl, Description string
	InstallFolder          string
}

func main() {

	items = []*Mod{
		&Mod{1, "My fantastic mod", "2.4", "Johny123", "http://magikinfo.pl"},
		&Mod{2, "Second super immersive", "1.3.55beta", "imbec", "https://www.magikinfo.cz"},
		&Mod{3, "Third super immersive", "1.4.55beta", "blabla", "https://www.magikinfo.pl"},
	}
	installed = []*Mod{
		&Mod{1, "Mod 1", "1", "auth", "http://localhost"},
		&Mod{2, "Mod 2", "1.5", "auth", "http://localhost"},
	}
	mp = &Modpack{
		Name:          "MyPack",
		Version:       "1.0",
		Author:        "ER",
		InstallFolder: "MyPack",
		SourceUrl:     "https://somesite.com/this/and/that",
		Description:   "This modpack bring some strange shit to your game. We will automaticaly break your game together. Common CTD's, inconsistencies, etc. expected.",
	}

	ui.Main(mainwindow)
}

func mainwindow() {
	w = ui.NewWindow("Windpeak", 640, 480, true)
	w.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		w.Destroy()
		return true
	})

	tab := ui.NewTab()
	w.SetChild(tab)
	w.SetMargined(true)

	tab.Append("Modpack", makeSettingPage(mp))
	tab.SetMargined(0, true)
	tab.Append("Missing mods", makeModsPage(items))
	tab.SetMargined(1, true)
	tab.Append("Already downloaded mods", makeModsPage(installed))
	tab.SetMargined(2, true)

	w.Show()
}

func makeSettingPage(mp *Modpack) ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	vbox.Append(ui.NewLabel(fmt.Sprintf("%s, %s", mp.Name, mp.Version)), false)
	vbox.Append(ui.NewLabel("Author: "+mp.Author), false)
	vbox.Append(ui.NewLabel("Homepage: "+mp.SourceUrl), false)
	vbox.Append(ui.NewLabel("Description:"), false)
	d := ui.NewMultilineEntry()
	d.SetText(mp.Description)
	d.SetReadOnly(true)
	vbox.Append(d, false)
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	hbox.Append(ui.NewButton("Select dir"), false)
	vbox.Append(hbox, false)

	return vbox
}

func makeModsPage(items []*Mod) ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	// keeps track of mod ID's and possition in vbox, it will change while deleting
	mIDs := []int{}

	for _, m := range items {
		mIDs = append(mIDs, m.ID)
		vbox.Append(makeModLine(m, vbox, &mIDs), false)
	}

	return vbox
}

func makeModLine(m *Mod, parent *ui.Box, mIDs *[]int) ui.Control {
	var i int
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	bUrl := ui.NewButton("Download")
	bUrl.OnClicked(func(f *ui.Button) {
		browser.OpenURL(m.URL)
		// deleting is nonsense here, must check file existence
		// but logic is good, reuse it!
		i, mIDs = remove(mIDs, m.ID) // removing current item from list, so indexes will change
		fmt.Println(i)
		parent.Delete(i)
	})

	hbox.Append(bUrl, false)
	hbox.Append(ui.NewLabel(m.Name), false)
	hbox.Append(ui.NewLabel(m.Version), false)
	hbox.Append(ui.NewLabel(m.Author), false)

	return hbox
}

// remove removes item from list and it's index
func remove(lp *[]int, item int) (int, *[]int) {
	l := *lp
	for i, other := range l {
		if other == item {
			ret := append(l[:i], l[i+1:]...)

			return i, &ret
		}
	}
	return 0, lp
}
