package main

import (
	"os"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		captureLogin("/etc/jacksquat.conf", 24*time.Hour, func() bool { return true })
	}
}
