package ekikara

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/mahonia"
	"net/http"
	"strconv"
)

type Schedule struct {
	Hour int64
	Min  int64
	To   string
}

type Ekikara struct {
	Station   string
	Direction string
}

func (e *Ekikara) getPage() *http.Response {
	resp, _ := http.Get("http://ekikara.jp/newdata/ekijikoku/" +
		e.Station +
		"/" +
		e.Direction +
		".htm")
	return resp
}

func (e *Ekikara) dump(t *html.Tokenizer) []Schedule {
	schedules := []Schedule{}
	var hour *int64 = nil
	var to *string = nil
	for {

		tt := t.Next()
		switch tt {
		case html.TextToken:
		case html.StartTagToken:
			//name, _ := t.TagName()
			more := true
			var attrs = map[string]string{}
			for more == true {
				key, val, m := t.TagAttr()
				more = m
				attrs[string(key)] = string(val)
			}
			if attrs["class"] == "textBold" {
				tt := t.Next()
				if tt == html.TextToken {
					hourText := string(t.Text())
					h, _ := strconv.ParseInt(hourText, 10, 64)
					hour = &h
				} else {
					t.Next()
					minString := string(t.Text())
					m, _ := strconv.ParseInt(minString, 10, 64)
					schedules = append(schedules, Schedule{
						Hour: *hour,
						Min:  m,
						To:   *to,
					})

				}
			}
			if attrs["class"] == "s" {
				t.Next()
				d := mahonia.NewDecoder("shift_jis")
				toString := string(t.Text())
				to_ := d.ConvertString(toString)
				to = &to_
			}

		case html.EndTagToken:
			endTagName, _ := t.TagName()
			if "table" == string(endTagName) {
				return schedules
			}
		}
		if tt == html.TextToken {
		}
	}
}

func (e *Ekikara) getTimeTableElement() []Schedule {
	schedules := []Schedule{}
	t := html.NewTokenizer(e.getPage().Body)
	for {
		tt := t.Next()
		if tt == html.ErrorToken {
			return schedules
		}
		if tt != html.StartTagToken {
			continue
		}
		more := true
		var attrs = map[string]string{}
		for more == true {
			key, val, m := t.TagAttr()
			more = m
			attrs[string(key)] = string(val)
		}
		if attrs["align"] != "center" {
			continue
		}

		if attrs["class"] != "lowBg06" {
			continue
		}
		//fmt.Println(e.dump(t))
		schedules = append(schedules, e.dump(t)...)
	}
	return nil
}

func (e *Ekikara) GetSchedules() []Schedule {
	return e.getTimeTableElement()
}

func NewEkikara(Station string, Direction string) *Ekikara {
	e := new(Ekikara)
	e.Station = Station
	e.Direction = Direction
	return e
}

/*
func main() {
	e := NewEkikara("1310071", "down1_13101231")
	schedules := e.GetSchedules()
	fmt.Println(schedules)
}

*/
