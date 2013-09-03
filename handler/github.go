package handler

import (
	"bytes"
	"fmt"
	"github.com/yushi/gother/github"
	"html/template"
	"net/http"
)

func notificationTable(notifications []github.Notification) string {
	var html bytes.Buffer
	t, err := template.New("gh_notification").Parse(`
<table id="projects" style="font-size: x-small">
  <tr>
    <th>Repo</th>
    <th>User</th>
    <th>Type</th>
    <th>At</th>
  </tr>
  {{range $i, $v:= .}}
  <tr>
    <td>{{.Repo}}</td>
    <td>{{.User}}</td>
    <td>{{.Type}}</td>
    <td>{{.At}}</td>
  </tr>
  {{end}}
</table>
`)

	err = t.Execute(&html, notifications)
	if err != nil {
		fmt.Println(err)
	}

	return html.String()
}

type GithubHandler struct {
}

func (gh *GithubHandler) HandleNotification(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var apiurl *string

	if len(q["user"]) == 0 {
		http.Error(w, "user parameter not found", http.StatusBadRequest)
		return
	}
	if len(q["apiurl"]) > 0 {
		apiurl = &q["apiurl"][0]
	}

	notifications := github.GetNotifications(q["user"][0], apiurl)
	fmt.Fprintf(w, "%s", notificationTable(notifications))
}
