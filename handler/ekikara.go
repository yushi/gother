package handler

import (
	"bytes"
	"fmt"
	"github.com/yushi/gother/ekikara"
	"html/template"
	"net/http"
)

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
    <td>{{.Hour}}{{.Min}}</td>
    <td>{{.To}}</td>
  </tr>
  {{end}}
</table>
`)

	err = t.Execute(&html, schedules)
	if err != nil {
		fmt.Println(err)
	}

	return html.String()
}

type EkikaraHandler struct {
	Schedules *[]ekikara.Schedule
}

func (p *EkikaraHandler) HandleEkikara(w http.ResponseWriter, r *http.Request) {
	if p.Schedules == nil {
		e := ekikara.NewEkikara("1310071", "down1_13101231")
		schedules := e.GetSchedules()
		p.Schedules = &schedules
	}
	fmt.Println(p.Schedules)
	fmt.Fprintf(w, "%s", trainTable(*p.Schedules))
}
