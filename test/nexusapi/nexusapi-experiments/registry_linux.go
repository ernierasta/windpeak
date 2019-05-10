package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
)

func registerNxm(binaryPath string) {

	dd := "[Desktop Entry]\nName=nxm\n" +
		fmt.Sprintf("Exec=%s --download %%u\n", binaryPath) +
		"Type=Application\n" +
		"Terminal=false\n" +
		"MimeType=x-scheme-handler/nxm;\n"

	home, _ := homedir.Dir()
	ioutil.WriteFile(home+"/.local/share/applications/nxm.desktop", []byte(dd), 0644)

	if _, err := os.Stat(home + "/.local/share/applications/mimeapps.list"); err == nil {
		// TODO: append under [Default Applications] section (use ini library for that?)
		log.Fatal("appending to ~/.local/share/applications/mimeapps.list UNIMPLEMENTED")
	} else if os.IsNotExist(err) {
		ld := "[Default Applications]\n" +
			"x-scheme-handler/nxm=nxm.desktop\n"

		ioutil.WriteFile(home+"/.local/share/applications/mimeapps.list", []byte(ld), 0644)
	} else {
		log.Fatal(err)
	}

	cmd := exec.Command("update-desktop-database", home+"/.local/share/applications")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
