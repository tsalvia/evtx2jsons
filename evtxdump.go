package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"github.com/0xrawsec/golang-evtx/evtx"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %[1]s FILES...\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()

	for _, evtxFile := range flag.Args() {
		ef, err := evtx.New(evtxFile)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		for e := range ef.FastEvents() {
			fmt.Println(string(evtx.ToJSON(e)))
		}
	}
}
