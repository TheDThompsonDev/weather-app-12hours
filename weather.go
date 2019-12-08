package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//types

type temperature struct {
	Value float64
	Unit  string
}

type weather struct {
	Datetime, IconPhrase     string
	Temperature              temperature
	PrecipitationProbability int
}

//ACCUWEATHER API

func generateURL() string {
	key := //put YOUR weather API key here
	req, _ := http.NewRequest(
		"GET",
		"http://dataservice.accuweather.com/forecasts/v1/hourly/12hour/16516_PC", //this location 16516_PC is for Cordova,tn. find your location in the location section 
		nil)
	q := req.URL.Query()
	q.Add("apikey", key)
	q.Add("metric", "false")
	q.Add("details", "true")
	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

func formatHour(h weather) string {
	hour := h.Datetime[11:16]
	return fmt.Sprintf("%s \t %s \t %.1f%s \t %d%%\n",
		hour,
		h.IconPhrase,
		h.Temperature.Value,
		h.Temperature.Unit,
		h.PrecipitationProbability,
	)
}

//local

func handler(w http.ResponseWriter, r *http.Request) {
	url := generateURL()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("failed in part 2")
	}

	body, err := ioutil.ReadAll(res.Body)
	var hours []weather
	e := json.Unmarshal(body, &hours)
	if e != nil {
		log.Fatal(err)
		fmt.Printf("failed in part 3")
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Weather</h1><ul>")
	for _, hour := range hours {
		fmt.Fprintf(w, "<li>%s</li>", formatHour(hour))
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


//accuweather is what I used to get the weather data. so apikey must be from there. 
