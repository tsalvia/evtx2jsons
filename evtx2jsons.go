package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func outputJsonFiles(outputDir string, stats []EventStats) {
	if _, err := os.Stat(outputDir); err != nil {
		if err = os.MkdirAll(outputDir, 0777); err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	for _, s := range stats {
		// Example: output/Security_4624.json
		outputFile := outputDir + "/" + s.Channel + "_" + strconv.FormatInt(s.EventID, 10) + ".json"
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
		fmt.Fprintln(of, "\t{}\n]") // {} is terminator.
	}
}

func parseTargetIDsOption(targetIDsOption string) []int64 {
	targetIDs := []int64{}
	if targetIDsOption != "" {
		sliceIDs := strings.Split(targetIDsOption, ",")
		for _, strID := range sliceIDs {
			id, _ := strconv.ParseInt(strID, 10, 64)
			targetIDs = append(targetIDs, id)
		}
	}
	return targetIDs
}

func containsTargetEventID(currentEventID int64, targetIDs []int64) bool {
	for _, targetID := range targetIDs {
		if currentEventID == targetID {
			return true
		}
	}
	return false
}

func containsEvent(stats []EventStats, eventID int64) (bool, int) {
	for i, s := range stats {
		if s.EventID == eventID {
			return true, i
		}
	}
	return false, -1
}

func evtx2json(evtxFile string, targetIDs []int64) []EventStats {
	ef, err := evtx.New(evtxFile)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	stats := []EventStats{}

	for e := range ef.FastEvents() {
		if len(targetIDs) != 0 && !containsTargetEventID(e.EventID(), targetIDs) {
			continue
		}

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
	var (
		input string
		outputDir string
		ids string
	)

	// Setting Options
	flag.StringVar(&input, "i", "", "This option is a short version of \"--input\" option.")
	flag.StringVar(&input, "input", "", "This option is required.\nSpecifies the EVTX file you want to convert to JSON file.")
	flag.StringVar(&outputDir, "d", "", "This option is a short version of \"--directory\" option.")
	flag.StringVar(&outputDir, "directory", "output", "Specifies the destination directory for the converted files. \n")
	flag.StringVar(&ids, "ids", "", "Specifies the event ID you want to output JOSN files.\nUse \",\" to separate multiple IDs.\n(default All Event IDs)")

	// Setting Help
	flag.Usage = func() {
		filename := filepath.Base(os.Args[0])
		// Usage
		fmt.Fprintf(os.Stderr, "\n%[1]s\n", filename)
		fmt.Fprintf(os.Stderr, "\n  Parse the EVTX file and output it in JSON format.\n")
		// Options
		fmt.Fprintf(os.Stderr, "\nOptions\n\n")
		flag.PrintDefaults()
		// Examples
		fmt.Fprintf(os.Stderr, "\nExamples\n\n")
		fmt.Fprintf(os.Stderr, "  1. Basic Usage\n")
		fmt.Fprintf(os.Stderr, "\t$ %s -i Security.evtx\n", filename)
		fmt.Fprintf(os.Stderr, "  2. Specify the event IDs you want to output.\n")
		fmt.Fprintf(os.Stderr, "\t$ %s -i Security.evtx -ids 4624,4625,1102 \n", filename)
		fmt.Fprintf(os.Stderr, "  3. Specify the destination directory.\n")
		fmt.Fprintf(os.Stderr, "\t$ %s -i Security.evtx -d output/jsons \n", filename)
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(0)
	}

	flag.Parse()

	if input == "" {
		fmt.Fprintf(os.Stderr, "Error: No EVTX file specified. Use the \"--input\" or \"-i\" option.\n")
		fmt.Fprintf(os.Stderr, "       For more information, try to use the \"--help\" or \"-h\" option to show help.\n")
		os.Exit(0)
	}

	targetIDs := parseTargetIDsOption(ids)

	stats := evtx2json(input, targetIDs)
	outputJsonFiles(outputDir, stats)
	showStats(stats)
}