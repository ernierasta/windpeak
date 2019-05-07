package main

import (
	"github.com/ernierasta/windpeak/gui"
)

func main() {
	s := gui.Settings{
		Modpack: &gui.Modpack{DirPath: "./mytest", Name: "test", Version: "1.2", Author: "ernierasta",
			Homepage: "https://github.com/ernierasta/windpeak", Description: "Testing my UI",
			Steps: []*gui.Step{
				&gui.Step{ID: 1, Name: "Basic patches", Description: "This is my great Description, what do you think?\nThis should be on new line. What about this long text. Does it mean anything?",
					Stop: true, StopCommands: []string{"runme.exe", "killall.exe", "C:\\Test\\test.exe --super param"},
				},
				&gui.Step{ID: 2, Name: "Stability"},
				&gui.Step{ID: 3, Name: "Graphics",
					Mods: []*gui.Mod{
						&gui.Mod{ID: 1, Name: "Testing nice very immersive mod you need to have", Author: "ernierasta", HomepageURL: "https://www.nexusmods.com/oblivion/mods/49335"},
					},
				},
			},
		},
	}
	g := gui.NewCreatorGUI(&s)
	g.Show()
}
