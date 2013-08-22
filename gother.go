package main

import (
	"encoding/json"
	"fmt"
	"github.com/yushi/gother/statusboard"
	"github.com/yushi/gother/system"
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
			"Used":     new([]statusboard.DataPoint),
			"Inactive": new([]statusboard.DataPoint),
			"Free":     new([]statusboard.DataPoint),
		}

		for _, memstat := range memstats {
			for memtype, datapoint := range datapoints {
				var val int64
				switch memtype {
				case "Used":
					val = memstat.meminfo.Wired + memstat.meminfo.Active
				case "Inactive":
					val = memstat.meminfo.Inactive
				case "Free":
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
			var color string
			switch memtype {
			case "Used":
				color = "Red"
			case "Inactive":
				color = "Blue"
			case "Free":
				color = "Green"
			}
			graph_entries = append(graph_entries,
				statusboard.GraphEntry{
					Title:      memtype,
					Color:      color,
					Datapoints: *datapoint,
				},
			)
		}
		jsonobj := statusboard.GraphJSON{
			Graph: statusboard.GraphData{
				Title:         "Memory",
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
	http.HandleFunc("/proc/meminfo", get_proc_handler())
	http.ListenAndServe(":8080", nil)
}
