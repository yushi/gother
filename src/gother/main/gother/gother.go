package main

import (
	"encoding/json"
	"fmt"
	"gother/statusboard"
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
		int_val = int_val * 4096        // page to Byte
		int_val = int_val / 1024 / 1024 // to MByte
		switch rows[0] {
		case "Pages free":
			m.free = int_val
		case "Pages speculative":
			m.free += int_val
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

func get_proc_handler() func(w http.ResponseWriter, r *http.Request) {
	type MemStat struct {
		label   string
		meminfo *MemInfo
	}
	memstats := make([]MemStat, 0)

	return func(w http.ResponseWriter, r *http.Request) {

		m := MemStat{
			label:   time.Now().Format("15:04:05"),
			meminfo: vm_stat(),
		}

		memstats = append(memstats, m)

		datapoints := map[string]*[]statusboard.DataPoint{
			"MemWired":    new([]statusboard.DataPoint),
			"MemActive":   new([]statusboard.DataPoint),
			"MemInactive": new([]statusboard.DataPoint),
			"MemFree":     new([]statusboard.DataPoint),
		}

		for _, memstat := range memstats {
			for memtype, datapoint := range datapoints {
				var val int
				switch memtype {
				case "MemWired":
					val = memstat.meminfo.wired
				case "MemActive":
					val = memstat.meminfo.active
				case "MemInactive":
					val = memstat.meminfo.inactive
				case "MemFree":
					val = memstat.meminfo.free
				}
				*datapoint = append(*datapoint,
					statusboard.DataPoint{
						Title: memstat.label,
						Value: val,
					})
			}
		}

		graph_entries := make([]statusboard.GraphEntry, 0)
		for memtype, datapoint := range datapoints {
			graph_entries = append(graph_entries,
				statusboard.GraphEntry{
					Title:      memtype,
					Datapoints: *datapoint,
				},
			)
		}
		jsonobj := statusboard.GraphJSON{
			Graph: statusboard.GraphData{
				Title:         "SystemInfo",
				Datasequences: graph_entries,
				Total:         false,
				Type:          "line",
			},
		}

		b, _ := json.Marshal(jsonobj)
		fmt.Fprintf(w, "%s", b)
	}
}

func main() {
	http.HandleFunc("/hello", hello_handler)
	http.HandleFunc("/proc", get_proc_handler())
	http.ListenAndServe(":8080", nil)
}
