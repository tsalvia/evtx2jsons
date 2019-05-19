package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"github.com/0xrawsec/golang-evtx/evtx"
)

type EventStats struct {
	Channel string
	EventID int64
	Count uint
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path) - len(filepath.Ext(path))])
}

func showStats(stats []EventStats) {
	fmt.Println("Channel, EventID, Count")
	for _, s := range stats {
		fmt.Printf("%s,\t%d,\t%d\n", s.Channel, s.EventID, s.Count)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %[1]s FILES...\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	stats := []EventStats{}
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
		defer of.Close()

		// Begin JSON
		fmt.Fprintln(of, "[")
		for e := range ef.FastEvents() {
			contains := false
			for i, s := range stats {
				if s.EventID == e.EventID() {
					stats[i].Count++
					contains = true
					break
				}
			}
			if !contains {
				newStats := EventStats{e.Channel(), e.EventID(), 1}
				stats = append(stats, newStats)
			}
			fmt.Fprintf(of, "\t%s,\n", string(evtx.ToJSON(e)))
		}
		fmt.Fprintf(of, "\t{}\n]")
		// End JSON
	}
	showStats(stats)
}
