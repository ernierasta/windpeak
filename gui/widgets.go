package gui

import (
	"github.com/andlabs/ui"
)

// here are custom widgest or interfaces

// DynamicBox allows to remove child and keep good order,
// so you will always remove correct entry.
// DynamicBox is always padded by default.
type DynamicBox struct {
	ids []int
	*ui.Box
}

func NewDynamicVbox() *DynamicBox {
	ids := []int{}
	b := DynamicBox{ids, ui.NewVerticalBox()}
	b.SetPadded(true)
	return &b
}

func (db *DynamicBox) Append(child ui.Control, id int) {
	db.ids = append(db.ids, id)
	db.Box.Append(child, false)
}

func (db *DynamicBox) Delete(id int) {
	var i int
	// removing current item from list, so indexes will change
	i, db.ids = RemoveSliceInt(db.ids, id)
	db.Box.Delete(i)
}

func (db *DynamicBox) DeleteAll() {
	for i := len(db.ids) - 1; i >= 0; i-- {
		db.Box.Delete(i)
	}
	db.ids = []int{}
}

type tab struct {
	id    int
	name  string
	child ui.Control
}

type DynamicTab struct {
	tabs []tab
	*ui.Tab
}

func NewDynamicTab() *DynamicTab {
	tabs := []tab{}
	t := DynamicTab{tabs, ui.NewTab()}
	return &t
}

func (dt *DynamicTab) Append(name string, child ui.Control, id int) {
	// remove all existing tabs (in reverse order is needed)
	for i := len(dt.tabs) - 1; i >= 0; i-- {
		dt.Tab.Delete(i)
	}
	// append new tab
	dt.Tab.Append(name, child)
	// reinsert rest of tabs
	for i := range dt.tabs {
		dt.Tab.InsertAt(dt.tabs[i].name, i, dt.tabs[i].child)
	}
	// store tab info internally
	dt.tabs = append(dt.tabs, tab{id, name, child})
}

//TODO: allow to keep tab last and keep desired order even if there is not enought tabs ATM

func (dt *DynamicTab) InsertAt(name string, pos int, child ui.Control, id int) {
	for i := len(dt.tabs) - 1; i >= 0; i-- {
		dt.Tab.Delete(i)
	}
	dt.Tab.Append(name, child)
	for i := range dt.tabs {
		if i < pos {
			dt.Tab.InsertAt(dt.tabs[i].name, i, dt.tabs[i].child)
		} else if i >= pos {
			dt.Tab.InsertAt(dt.tabs[i].name, i+1, dt.tabs[i].child)
		}
	}
	dt.tabs = InsertAtSliceTab(dt.tabs, tab{id, name, child}, pos)
	//dt.Tab.InsertAt(name, order, child)
}

func (dt *DynamicTab) SetMargined(id int, margined bool) {
	i := IndexSliceTab(dt.tabs, id)
	if i != -1 {
		dt.Tab.SetMargined(i, margined)
	}
}

func (dt *DynamicTab) Delete(id int) {
	var i int
	// removing current item from list, so indexes will change
	i, dt.tabs = RemoveSliceTab(dt.tabs, id)
	dt.Tab.Delete(i)
}

// ComboButton is special type of button allowing you to
// Show/Hide combo. You need to provide your combo
// (allows to place it where you need it).
type ComboButton struct {
	options   []string
	selected  int
	onClicked func(*ui.Button)
	cbox      *ui.Combobox
	*ui.Button
}

func NewComboButton(buttonTitle string, options []string, combo *ui.Combobox) *ComboButton {
	cbtn := ComboButton{}
	cbtn.Button = ui.NewButton(buttonTitle)
	cbtn.cbox = combo
	cbtn.options = options
	cbtn.selected = -1 // to distinguish if option was selected
	for i := range options {
		cbtn.cbox.Append(options[i])
	}
	cbtn.cbox.OnSelected(func(c *ui.Combobox) {
		cbtn.cbox.Hide()
		cbtn.selected = cbtn.cbox.Selected()
		cbtn.onClicked(cbtn.Button)
	})
	cbtn.cbox.Hide()

	return &cbtn
}

// OnClicked will set OnSelected function for combo and
// set hiding/showing combo for button.
func (cbtn *ComboButton) OnClicked(f func(b *ui.Button)) {
	cbtn.onClicked = f
	cbtn.Button.OnClicked(func(b *ui.Button) {
		if !cbtn.cbox.Visible() {
			cbtn.cbox.SetSelected(-1) // reset previous selection
			cbtn.cbox.Show()
		} else {
			cbtn.cbox.Hide()
		}
	})

}

func (cbtn *ComboButton) Selected() int {
	return cbtn.selected
}

func RemoveSliceInt(sl []int, item int) (int, []int) {
	for i, other := range sl {
		if other == item {
			ret := append(sl[:i], sl[i+1:]...)

			return i, ret
		}
	}
	return 0, sl
}

func RemoveSliceTab(sl []tab, id int) (int, []tab) {
	for i, other := range sl {
		if other.id == id {
			ret := append(sl[:i], sl[i+1:]...)

			return i, ret
		}
	}
	return 0, sl
}

func IndexSliceTab(sl []tab, id int) int {
	for i, other := range sl {
		if other.id == id {
			return i
		}
	}
	return -1
}

func InsertAtSliceInt(sl []int, item, index int) []int {
	sl = append(sl, 0)
	copy(sl[index+1:], sl[index:])
	sl[index] = item
	return sl
}

func InsertAtSliceTab(sl []tab, item tab, index int) []tab {
	sl = append(sl, tab{})
	copy(sl[index+1:], sl[index:])
	sl[index] = item
	return sl
}
