package clr

import "fmt"

const (
	Blk = "\u001b[1;30m"
	Red = "\u001b[1;31m"
	Grn = "\u001b[1;32m"
	Yel = "\u001b[1;33m"
	Blu = "\u001b[1;34m"
	Cyn = "\u001b[1;36m"
	Off = "\u001b[0m"
)

func Print(clr string, msg string) {
	fmt.Print(clr, msg, Off)
}

func Printf(clr string, msg string, v ...interface{}) {
	fmt.Print(clr)
	fmt.Printf(msg, v...)
	fmt.Print(Off)
}
