package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type MemInfo struct {
	free     int
	active   int
	wired    int
	inactive int
}

func vm_stat() *MemInfo {
	m := new(MemInfo)
	out, err := exec.Command("vm_stat").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	lines = lines[1 : len(lines)-2] // exclude header/footer
	for _, line := range lines {
		rows := strings.Split(line, ":")
		val := rows[1]
		val = strings.TrimLeft(val, " ")
		val = strings.TrimRight(val, ".")
		int_val, _ := strconv.Atoi(val)
		switch rows[0] {
		case "Pages free":
			m.free = int_val
		case "Pages active":
			m.active = int_val
		case "Pages inactive":
			m.inactive = int_val
		case "Pages wired down":
			m.wired = int_val
		}
	}
	return m
}

func hello_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go!")
}

type StatusBoardJSON struct {
	Graph StatusBoardJSONGraph `json:"graph"`
}

type StatusBoardJSONGraph struct {
	Title         string                      `json:"title"`
	Datasequences []StatusBoardJSONGraphEntry `json:"datasequences"`
}

type StatusBoardJSONGraphEntry struct {
	Title      string                     `json:"title"`
	Datapoints []StatusBoardJSONDataPoint `json:"datapoints"`
}

type StatusBoardJSONDataPoint struct {
	Title string `json:"title"`
	Value int    `json:"value"`
}

func proc_handler(w http.ResponseWriter, r *http.Request) {
	m := vm_stat()
	memfree_datapoints := []StatusBoardJSONDataPoint{
		StatusBoardJSONDataPoint{
			Title: "20130820",
			Value: m.free,
		},
		StatusBoardJSONDataPoint{
			Title: "20130821",
			Value: m.free,
		},
		StatusBoardJSONDataPoint{
			Title: "20130822",
			Value: m.free,
		},
	}
	memactive_datapoints := []StatusBoardJSONDataPoint{
		StatusBoardJSONDataPoint{
			Title: "20130820",
			Value: m.active,
		},
		StatusBoardJSONDataPoint{
			Title: "20130821",
			Value: m.active,
		},
		StatusBoardJSONDataPoint{
			Title: "20130822",
			Value: m.active,
		},
	}

	graph_entries := []StatusBoardJSONGraphEntry{
		StatusBoardJSONGraphEntry{
			Title:      "MemFree",
			Datapoints: memfree_datapoints,
		},
		StatusBoardJSONGraphEntry{
			Title:      "MemActive",
			Datapoints: memactive_datapoints,
		},
	}

	jsonobj := StatusBoardJSON{
		Graph: StatusBoardJSONGraph{
			Title:         "SystemInfo",
			Datasequences: graph_entries,
		},
	}

	fmt.Println(jsonobj)
	b, err := json.Marshal(jsonobj)
	fmt.Println(err)
	fmt.Println(time.Now())
	fmt.Fprintf(w, "%s", b)
}

func main() {
	http.HandleFunc("/hello", hello_handler)
	http.HandleFunc("/proc", proc_handler)
	http.ListenAndServe(":8080", nil)
}
