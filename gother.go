package main

import (
	"flag"
	"fmt"
	"github.com/yushi/gother/statusboard"
	"github.com/yushi/gother/system"
	"log"
	"net/http"
	"time"
)

func hello_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go!")
}

func getTimeStr() string {
	// for test
	//return time.Now().Format("15:04:05")

	return time.Now().Format("15:04")
}

func updateStats(stats []system.StatHistory) []system.StatHistory {
	now := getTimeStr()
	if len(stats) == 0 || now != stats[len(stats)-1].Time {
		stats = append(stats,
			system.StatHistory{
				Time: getTimeStr(),
				Stat: system.GetStat(),
			})
	}
	if len(stats) > 1440 {
		stats = stats[0:1440]
	}
	return stats
}

func getProcLoadHandler() func(w http.ResponseWriter, r *http.Request) {
	stats := make([]system.StatHistory, 0)

	return func(w http.ResponseWriter, r *http.Request) {
		stats = updateStats(stats)
		fmt.Fprintf(w, "%s", statusboard.LoadavgGraph(stats))
	}
}

func getProcMemHandler() func(w http.ResponseWriter, r *http.Request) {
	stats := make([]system.StatHistory, 0)

	return func(w http.ResponseWriter, r *http.Request) {
		stats = updateStats(stats)
		fmt.Fprintf(w, "%s", statusboard.MemoryGraph(stats))
	}
}

func main() {
	port := flag.Int("p", 9090, "listen port")
	flag.Parse()

	http.HandleFunc("/hello", hello_handler)
	http.HandleFunc("/proc/mem", getProcMemHandler())
	http.HandleFunc("/proc/load", getProcLoadHandler())

	log.Printf("About to listen on %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
