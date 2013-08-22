package statusboard

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
