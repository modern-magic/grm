package main

import (
	"fmt"
	"nrm/internal"
	"os"
)

var (
	version = "V0.0.0"
	date    = "0000-00-00T00:00:00Z"
)

func main() {
	/*
	 we can get os args form os.Args it will array .
	*/
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Printf("gonrm %s (%s)\n", version, date)
		os.Exit(0)
		return
	}
	internal.ShowList()
}
