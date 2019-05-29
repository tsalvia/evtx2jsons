package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"path/filepath"
	"github.com/0xrawsec/golang-evtx/evtx"
)

type EventStats struct {
	Channel string
	EventID int64
	Count uint
	EvtxJsons []string
}

func showStats(stats []EventStats) {
	fmt.Println("Channel, EventID, Count")
	for _, s := range stats {
		fmt.Printf("%s,\t%d,\t%d\n", s.Channel, s.EventID, s.Count)
	}
}

func outputJsonFiles(stats []EventStats) {
	for _, s := range stats {
		outputFile := s.Channel + "_" + strconv.FormatInt(s.EventID, 10) + ".json"
		of, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer of.Close()

		fmt.Fprintln(of, "[")
		for _, evtxJson := range s.EvtxJsons {
			fmt.Fprintf(of, "\t%s,\n", evtxJson)
		}
		fmt.Fprintln(of, "\t{}\n]")
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

	stats := []EventStats{}

	for e := range ef.FastEvents() {
		contains, num := containsEvent(stats, e.EventID())
		evtxJson := string(evtx.ToJSON(e))
		if !contains {
			newStats := EventStats{e.Channel(), e.EventID(), 1, []string{evtxJson}}
			stats = append(stats, newStats)
		} else {
			stats[num].Count++
			stats[num].EvtxJsons = append(stats[num].EvtxJsons, evtxJson)
		}
	}

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
	outputJsonFiles(stats)
	showStats(stats)
}