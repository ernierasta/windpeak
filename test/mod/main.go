package main

import (
	"fmt"
	"time"

	"github.com/ernierasta/windpeak/mod"
)

func main() {
	dd := "/home/ernie/Temp/delme"

	gd := "/home/ernie/Temp"
	m := mod.NewWreyBashMod(dd, gd)
	t := time.Now()
	m.CreateMeta(false, true)
	fmt.Printf("xxhash: %f\n", time.Since(t).Seconds())
	t = time.Now()
	m.CreateMeta(false, false)
	fmt.Printf("md5hash: %f\n", time.Since(t).Seconds())
	t = time.Now()
	m.CreateMeta(true, true)
	fmt.Printf("no hash: %f\n", time.Since(t).Seconds())
}
