package github

import (
	"github.com/google/go-github/github"
	"net/url"
	"strings"
)

type Notification struct {
	User string
	Repo string
	Type string
	At   string
}

func APIURL(urlstr string) *url.URL {
	parsed, _ := url.Parse(urlstr)
	return parsed
}

func GetNotifications(username string, apiUrl *string) []Notification {
	client := github.NewClient(nil)
	if apiUrl != nil {
		client.BaseURL = APIURL(*apiUrl)
	}
	ev, _, _ := client.Activity.ListEventsRecievedByUser(
		username,
		true,
		&github.ListOptions{Page: 1})

	notifications := make([]Notification, 0)
	for _, event := range ev {
		n := Notification{
			User: *event.Actor.Login,
			Repo: *event.Repo.Name,
			Type: strings.TrimSuffix(*event.Type, "Event"),
			At:   event.CreatedAt.Local().Format("1/2 15:04"),
		}
		notifications = append(notifications, n)
	}
	return notifications
}
