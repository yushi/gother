package handler

import (
	"bytes"
	"fmt"
	"github.com/yushi/gother/ekikara"
	"html/template"
	"net/http"
	"time"
)

type ReadableSchedule struct {
	Time string
	To   string
}

func trainTable(schedules []ekikara.Schedule) string {
	var html bytes.Buffer
	t, err := template.New("ekikara").Parse(`
<table id="projects">
  <tr>
    <th>Time</th>
    <th>Dst</th>
  </tr>
  {{range $i, $v:= .}}
  <tr>
    <td>{{.Time}}</td>
    <td>{{.To}}</td>
  </tr>
  {{end}}
</table>
`)

	now := time.Now()

	vars := []ReadableSchedule{}
	for _, s := range schedules {
		if s.Hour < int64(now.Hour()) {
			continue
		}
		if s.Hour == int64(now.Hour()) && s.Min < int64(now.Minute()) {
			continue
		}

		vars = append(vars, ReadableSchedule{
			Time: fmt.Sprintf("%02d:%02d", s.Hour, s.Min),
			To:   s.To,
		})
		if len(vars) > 9 {
			break
		}
	}
	err = t.Execute(&html, vars)
	if err != nil {
		fmt.Println(err)
	}

	return html.String()
}

type EkikaraHandler struct {
	Schedules map[string][]ekikara.Schedule
}

func (p *EkikaraHandler) Init() {
	p.Schedules = make(map[string][]ekikara.Schedule)
}

func (p *EkikaraHandler) HandleEkikara(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	dir := q["dir"][0]
	file := q["file"][0]
	key := dir + file

	if _, ok := p.Schedules[key]; !ok {
		e := ekikara.NewEkikara(dir, file)
		schedules := e.GetSchedules()
		p.Schedules[key] = schedules
	}

	fmt.Fprintf(w, "%s", trainTable(p.Schedules[key]))
}
