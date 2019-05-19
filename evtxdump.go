package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"github.com/0xrawsec/golang-evtx/evtx"
)

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path) - len(filepath.Ext(path))])
}


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

		outputFile := getFileNameWithoutExt(evtxFile) + ".json"
		of, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Fprintln(of, "[")

		for e := range ef.FastEvents() {
			fmt.Fprintf(of, "\t%s,\n", string(evtx.ToJSON(e)))
		}

		fmt.Fprintln(of, "\t{}\n]")
		of.Close()

	}
}
