package main

import (
	"encoding/json"
	"fmt"
	"github.com/yushi/gother/statusboard"
	"github.com/yushi/gother/system"
	"net/http"
	"sort"
	"time"
)

func hello_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go!")
}

func get_proc_cpu_handler() func(w http.ResponseWriter, r *http.Request) {
	type StatHistory struct {
		label string
		stat  system.SystemStat
	}
	stats := make(map[string]*system.SystemStat)

	return func(w http.ResponseWriter, r *http.Request) {
		stats[time.Now().Format("15:04:05")] = system.GetSystemStat()
		datapoints := make(map[string][]statusboard.DataPoint)
		for _, val := range []string{"sys", "user", "idle"} {
			keys := make([]string, 0)
			for k, _ := range stats {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				var v float64
				s := stats[k]
				switch val {
				case "sys":
					v = s.Cpu.Sys
				case "user":
					v = s.Cpu.User
				case "idle":
					v = s.Cpu.Idle
				}
				datapoints[val] = append(
					datapoints[val],
					statusboard.DataPoint{
						Title: k,
						Value: v,
					})
			}
		}

		graph_entries := make([]statusboard.GraphEntry, 0)
		for cpu_type, datapoint := range datapoints {
			var color string
			switch cpu_type {
			case "Used":
				color = "Red"
			case "Inactive":
				color = "Blue"
			case "Free":
				color = "Green"
			}
			graph_entries = append(graph_entries,
				statusboard.GraphEntry{
					Title:      cpu_type,
					Color:      color,
					Datapoints: datapoint,
				},
			)
		}
		jsonobj := statusboard.GraphJSON{
			Graph: statusboard.GraphData{
				Title:         "CPU Usage",
				Datasequences: graph_entries,
				Total:         false,
				Type:          "line",
			},
		}

		b, _ := json.Marshal(jsonobj)
		fmt.Fprintf(w, "%s", b)

	}
}

func get_proc_mem_handler() func(w http.ResponseWriter, r *http.Request) {
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
				var val float64
				switch memtype {
				case "Used":
					val = float64(memstat.meminfo.Wired) +
						float64(memstat.meminfo.Active)
				case "Inactive":
					val = float64(memstat.meminfo.Inactive)
				case "Free":
					val = float64(memstat.meminfo.Free)
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
	http.HandleFunc("/proc/mem", get_proc_mem_handler())
	http.HandleFunc("/proc/cpu", get_proc_cpu_handler())
	http.ListenAndServe(":8080", nil)
}
