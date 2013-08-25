package handler

import (
	"fmt"
	"github.com/yushi/gother/statusboard"
	"github.com/yushi/gother/system"
	"net/http"
	"time"
)

type ProcHandler struct {
	Stats []system.StatHistory
}

func (p *ProcHandler) HandleLoadavg(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", statusboard.LoadavgGraph(p.Stats))
}

func (p *ProcHandler) HandleMemory(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", statusboard.MemoryGraph(p.Stats))
}

func getTimeStr() string {
	// for test
	//return time.Now().Format("15:04:05")

	return time.Now().Format("15:04")
}

func (p *ProcHandler) Update() {
	now := getTimeStr()
	if len(p.Stats) == 0 || now != p.Stats[len(p.Stats)-1].Time {
		p.Stats = append(p.Stats,
			system.StatHistory{
				Time: getTimeStr(),
				Stat: system.GetStat(),
			})
	}
	if len(p.Stats) > 1440 {
		p.Stats = p.Stats[0:1440]
	}
}

func (p *ProcHandler) UpdatePeriodically() {
	for {
		p.Update()
		time.Sleep(30 * time.Second)
	}
}
func (p *ProcHandler) Start() {
	go p.UpdatePeriodically()
}
