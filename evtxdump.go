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

func containsEvent(stats []EventStats, eventID int64) (bool, int) {
	for i, s := range stats {
		if s.EventID == eventID {
			return true, i
		}
	}
	return false, -1
}

func evtx2json(evtxFile string) []EventStats {
	ef, err := evtx.New(evtxFile)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	outputFile := getFileNameWithoutExt(evtxFile) + ".json"
	of, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer of.Close()

	stats := []EventStats{}

	// Begin JSON
	fmt.Fprintln(of, "[")
	for e := range ef.FastEvents() {
		contains, num := containsEvent(stats, e.EventID())
		if !contains {
			newStats := EventStats{e.Channel(), e.EventID(), 1}
			stats = append(stats, newStats)
		} else {
			stats[num].Count++
		}
		// Append Event to JSON
		fmt.Fprintf(of, "\t%s,\n", string(evtx.ToJSON(e)))
	}
	fmt.Fprintf(of, "\t{}\n]")
	// End JSON

	return stats
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
		newStats := evtx2json(evtxFile)
		stats = append(stats, newStats...)
	}
	showStats(stats)
}