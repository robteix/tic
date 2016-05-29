package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	meaningfulDates = []MeaningfulDate{}
	ticTmpl         *template.Template
	port            = ":8080"
)

// getTic handles GET /
func getTic(w http.ResponseWriter, r *http.Request) {

	type line struct {
		Desc string
		Tic  string
	}
	data := make([]line, len(meaningfulDates))
	for k, v := range meaningfulDates {
		data[k] = line{
			Desc: v.Desc,
			Tic:  formatDateDiff(time.Now(), v.Date),
		}
	}
	ticTmpl.Execute(w, data)
}

// formatDateDiff formats the difference between two
// dates as years, months, and days.
func formatDateDiff(a, b time.Time) string {
	year, month, day, _, _, _ := dateDiff(a, b)
	parts := []string{}
	if year > 0 {
		parts = append(parts, fmt.Sprintf("%dy ", year))
	}
	if month > 0 {
		parts = append(parts, fmt.Sprintf("%dm ", month))
	}
	if day > 0 {
		parts = append(parts, fmt.Sprintf("%dd ", day))
	}

	return strings.Join(parts, " ")
}

// dateDiff calculates the difference between two dates as
// years, months, days, hours, minutes, and seconds.
func dateDiff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// MeaningfulDate represents a meaningful date
// that willl be displayed
type MeaningfulDate struct {
	Desc string
	Date time.Time
}

func main() {

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	torontoLoc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		log.Fatalf("error loading timezone info: %v", err)
	}

	ticTmpl, err = template.New("tmpl").Parse(tmpl)
	if err != nil {
		log.Fatalf("error loading template: %v", err)
	}
	meaningfulDates = []MeaningfulDate{
		MeaningfulDate{
			Desc: "Beto's TIC",
			Date: time.Date(2016, 4, 7, 0, 0, 0, 0, torontoLoc),
		},
		MeaningfulDate{
			Desc: "Girls' TIC",
			Date: time.Date(2016, 4, 27, 0, 0, 0, 0, torontoLoc),
		},
		MeaningfulDate{
			Desc: "In job",
			Date: time.Date(2016, 4, 25, 0, 0, 0, 0, torontoLoc),
		},
	}

	http.HandleFunc("/", getTic)
	http.ListenAndServe(port, nil)
}

var tmpl = `{{define "tmpl"}}
<!DOCTYPE html>
<head>
<title>TIC</title>
<style>
h1 {
    color: red;
}
</style>
</head>
<body>
<center>
{{range .}}
<h2>{{.Desc}}</h2>
<h1>{{.Tic}}</h1>
{{end}}
</center>
</body>
</html>
{{end}}`
