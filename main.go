package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Parameters struct {
	ChannelName string
	Precision   string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("*.gohtml"))
}

func main() {
	http.HandleFunc("/", math)
	http.ListenAndServe(":8080", nil)
}

func math(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel")
	//p := Parameters{}
	//p.ChannelName = "Leprechaunkoala"
	//p.Precision = "precision=4"

	resp, err := http.Get("https://decapi.me/twitch/uptime/" + channel + "?precision=4")

	if err != nil {
		fmt.Println("Request Failed: $s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	BodyString := string(body)
	// fmt.Println(strings.Split(bodyString, ","))

	// exString := "2 hours, 20 minutes, 30 seconds"
	time := strings.Split(BodyString, ",")

	type Time struct {
		Hours   string
		Minutes string
		Seconds string
	}

	t := Time{}

	hours, err := regexp.MatchString("hours?", time[0])
	minutes, err := regexp.MatchString("minutes?", time[0])
	seconds, err := regexp.MatchString("seconds?", time[0])

	regLetters := regexp.MustCompile(`[^0-9]`)

	if hours == true {
		t.Hours = regLetters.ReplaceAllString(time[0], "")   // strings.Replace(time[0], "hours", "", -1)
		t.Minutes = regLetters.ReplaceAllString(time[1], "") // strings.Replace(time[1], "minutes", "", -1)
		t.Seconds = regLetters.ReplaceAllString(time[2], "") // strings.Replace(time[2], "seconds", "", -1)
	} else if minutes == true {
		t.Hours = "0"
		t.Minutes = regLetters.ReplaceAllString(time[0], "") // strings.Replace(time[0], "minutes", "", -1)
		t.Seconds = regLetters.ReplaceAllString(time[1], "") // strings.Replace(time[1], "seconds", "", -1)
	} else if seconds == true {
		t.Hours = "0"
		t.Minutes = "0"
		t.Seconds = regLetters.ReplaceAllString(time[0], "") // strings.Replace(time[0], "seconds", "", -1)
	}

	hoursN, err := strconv.ParseFloat(t.Hours, 64)
	minutesN, err := strconv.ParseFloat(t.Minutes, 64)
	secondsN, err := strconv.ParseFloat(t.Seconds, 64)

	hydrate := hoursN*1000 + minutesN*16.66 + secondsN*0.27

	hydrateTotal := strconv.FormatFloat(hydrate, 'f', 2, 64)

	data := struct {
		Uptime  string
		Hydrate string
	}{
		Uptime:  BodyString,
		Hydrate: hydrateTotal,
	}

	tpl.ExecuteTemplate(w, "index.gohtml", data)
	// fmt.Println("Stream has been live for" + " " + bodyString + "," + " " + "You should have drunk" + " " + HydrateTotal + "ml of water so far")
}
