package statusboard

import (
	"encoding/json"
	"github.com/yushi/gother/system"
	"reflect"
	"sort"
)

type GraphJSON struct {
	Graph GraphData `json:"graph"`
}

type GraphData struct {
	Title         string       `json:"title"`
	Total         bool         `json:"total"`
	Type          string       `json:"type"`
	Datasequences []GraphEntry `json:"datasequences"`
}

type GraphEntry struct {
	Title      string      `json:"title"`
	Color      string      `json:"color",omitempty`
	Datapoints []DataPoint `json:"datapoints"`
}

type DataPoint struct {
	Title string  `json:"title"`
	Value float64 `json:"value"`
}

func LoadavgGraph(stats map[string]*system.SystemStat) []byte {
	datapoints := getDatapoints(stats, "Load", []string{"Load1", "Load5", "Load15"})

	color_map := map[string]string{
		"used":   "Red",
		"cached": "Blue",
		"free":   "Green",
	}

	graph_entries := getGraphEntries(datapoints, color_map)
	return getGraphJSON("Loadavg", graph_entries)
}

func MemoryGraph(stats map[string]*system.SystemStat) []byte {
	datapoints := getDatapoints(stats, "Mem", []string{"Used", "Cached", "Free"})
	color_map := map[string]string{
		"Load1":  "Red",
		"Load5":  "Blue",
		"Load15": "Green",
	}

	graph_entries := getGraphEntries(datapoints, color_map)
	return getGraphJSON("Memory", graph_entries)
}

func getDatapoints(stats map[string]*system.SystemStat, statField string, valueFields []string) map[string][]DataPoint {
	datapoints := make(map[string][]DataPoint)
	for _, val := range valueFields {
		keys := make([]string, 0)
		for k, _ := range stats {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			statValue := reflect.ValueOf(stats[k]).Elem()
			statFieldValue := statValue.FieldByName(statField)
			v := statFieldValue.Elem().FieldByName(val).Float()
			datapoints[val] = append(
				datapoints[val],
				DataPoint{
					Title: k,
					Value: v,
				})
		}
	}
	return datapoints
}

func getGraphEntries(datapoints map[string][]DataPoint, color_map map[string]string) []GraphEntry {

	graph_entries := make([]GraphEntry, 0)
	for k, datapoint := range datapoints {
		graph_entries = append(graph_entries,
			GraphEntry{
				Title:      k,
				Color:      color_map[k],
				Datapoints: datapoint,
			},
		)
	}
	return graph_entries
}

func getGraphJSON(title string, graph_entries []GraphEntry) []byte {
	jsonobj := GraphJSON{
		Graph: GraphData{
			Title:         title,
			Datasequences: graph_entries,
			Total:         false,
			Type:          "line",
		},
	}

	b, _ := json.Marshal(jsonobj)
	return b
}
