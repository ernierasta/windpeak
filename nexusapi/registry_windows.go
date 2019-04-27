package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

// register_nxm registers Windpeak as nxm handler
func register_nxm(path string) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\nxm\shell\open`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	if err := k.SetStringValue("command", fmt.Sprintf("\"%s\" \"%1\""), path); err != nil {
		log.Fatal(err)
	}
	if err := k.Close(); err != nil {
		log.Fatal(err)
	}
}
