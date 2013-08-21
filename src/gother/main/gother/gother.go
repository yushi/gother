package main

import (
	"encoding/json"
	"fmt"
	"gother/statusboard"
	"gother/system"
	"net/http"
	"time"
)

func hello_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go!")
}

func get_proc_handler() func(w http.ResponseWriter, r *http.Request) {
	type MemStat struct {
		label   string
		meminfo *system.MemInfo
	}
	memstats := make([]MemStat, 0)

	return func(w http.ResponseWriter, r *http.Request) {

		m := MemStat{
			label:   time.Now().Format("15:04:05"),
			meminfo: system.GetMemInfo(),
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
					val = memstat.meminfo.Wired
				case "MemActive":
					val = memstat.meminfo.Active
				case "MemInactive":
					val = memstat.meminfo.Inactive
				case "MemFree":
					val = memstat.meminfo.Free
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
