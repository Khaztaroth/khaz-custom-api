package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
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
	http.HandleFunc("/hydratebot", hydration)
	http.HandleFunc("/weather", weather)
	http.ListenAndServe(":8080", nil)
}

func hydration(w http.ResponseWriter, r *http.Request) {
	//URL querry for getting channel name
	channel := r.URL.Query().Get("channel")
	//simple := r.URL.Query().Get("no_message")

	//GET request to DecApi twitch channel uptime endpoint
	resp, err := http.Get("https://decapi.me/twitch/uptime/" + channel + "?precision=4")

	if err != nil {
		fmt.Println("Request Failed: $s", err)
	}
	//Storing body from DecApi request as variable
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	BodyString := string(body)

	// Example string for testing
	// exString := "2 hours, 20 minutes, 30 seconds"

	//Splitting string into hours, minutes, and seconds
	time := strings.Split(BodyString, ",")

	//Struct for each time measure and declaring a variable to use the info
	type Time struct {
		Hours   string
		Minutes string
		Seconds string
	}
	t := Time{}

	//Regex checks to see what kind of time measure we get
	hours, err := regexp.MatchString("hours?", time[0])
	minutes, err := regexp.MatchString("minutes?", time[0])
	seconds, err := regexp.MatchString("seconds?", time[0])

	//Regex setup to remove anything that's not a number later
	regLetters := regexp.MustCompile(`[^0-9]`)

	//Codeblock to check time measure structure and removing anything that's not a number. Stored as a value for later use
	if hours == true {
		t.Hours = regLetters.ReplaceAllString(time[0], "")
		t.Minutes = regLetters.ReplaceAllString(time[1], "")
		t.Seconds = regLetters.ReplaceAllString(time[2], "")
	} else if minutes == true {
		t.Hours = "0"
		t.Minutes = regLetters.ReplaceAllString(time[0], "")
		t.Seconds = regLetters.ReplaceAllString(time[1], "")
	} else if seconds == true {
		t.Hours = "0"
		t.Minutes = "0"
		t.Seconds = regLetters.ReplaceAllString(time[0], "")
	}
	//Converting string into float64 values to use them in math calculations
	hoursFloat, err := strconv.ParseFloat(t.Hours, 64)
	minutesFloat, err := strconv.ParseFloat(t.Minutes, 64)
	secondsFloat, err := strconv.ParseFloat(t.Seconds, 64)

	//Converting time into seconds
	Seconds := hoursFloat*3600 + minutesFloat*60 + secondsFloat

	//Calculating water intake based on each time measure
	hydrationML := Seconds * 0.213097
	hydrationL := hydrationML / 1000
	hydrationOz := hydrationML * 0.0338140227

	//Max water intake per day
	MaxWaterIntakeLitre := 3.7
	MaxWaterIntakeOunce := 130

	if hydrationL > float64(MaxWaterIntakeLitre) {
		hydrationL = 3.7
	}
	if hydrationOz > float64(MaxWaterIntakeOunce) {
		hydrationOz = 130.22
	}

	//Transforming float64 values to string to use in the final message
	volumeL := strconv.FormatFloat(hydrationL, 'f', 2, 64) + "L"
	volumeOz := strconv.FormatFloat(hydrationOz, 'f', 2, 64) + "fl oz"

	//Checking if request to DecApi returns the channel as 'offline', used to set a different message
	offline, err := regexp.MatchString("offline", time[0])

	//Struct to encode resulting message, this information is passed to hte gohtml template
	type Data struct {
		//Simple  string
		User    string
		Offline bool
		Uptime  string
		ML      string
		OZ      string
	}
	info := Data{channel, offline, BodyString, volumeL, volumeOz}
	//fmt.Println(info)

	tpl.ExecuteTemplate(w, "hydration.gohtml", info)
}

func weather(w http.ResponseWriter, r *http.Request) {
	//URL querry for getting the place
	place := r.URL.Query().Get("place")
	//simple := r.URL.Query().Get("no_message")

	//GET request to ScorpStuff weather API
	resp, err := http.Get("http://api.scorpstuff.com/weather.php?units=metric&city=" + url.QueryEscape(place))

	if err != nil {
		fmt.Println("Request Failed: $s", err)
	}
	//Storing body from DecApi request as variable
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	BodyString := string(body)

	CensoredData := strings.Replace(BodyString, "Chillán, CL", "Birbland", -1)

	type Data struct {
		BodyString   string
		CensoredData string
	}

	info := Data{BodyString, CensoredData}

	tpl.ExecuteTemplate(w, "weather.gohtml", info)

	// fmt.Println(CensoredData)
}
