package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

// registerNxm registers Windpeak as nxm handler
func registerNxm(path string) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\nxm\shell\open\command`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	if err := k.SetStringValue("", fmt.Sprintf("\"%s\" --download \"%%1\"", path)); err != nil {
		log.Fatal(err)
	}
	if err := k.Close(); err != nil {
		log.Fatal(err)
	}
}
