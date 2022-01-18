package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

type StatusData struct {
	WaterStatus struct {
		Water  int    `json:"water"`
		Status string `json:"statusWater"`
	}
	WindStatus struct {
		Wind   int    `json:"wind"`
		Status string `json:"statusWind"`
	}
}

func main() {
	go AutoReloadJSON()
	http.HandleFunc("/", AutoReloadWeb)
	fmt.Println("listening on PORT: ", "8080")
	http.ListenAndServe(":8080", nil)
}

func AutoReloadJSON() {
	for {
		min := 1
		max := 21

		water := rand.Intn(max-min) + 1
		wind := rand.Intn(max-min) + 1

		data := StatusData{}
		data.WaterStatus.Water = water
		switch {
		case water < 5:
			data.WaterStatus.Status = "Aman"
		case water > 5 && water <= 8:
			data.WaterStatus.Status = "Siaga"
		case water > 8:
			data.WaterStatus.Status = "Bahaya"
		default:
			data.WaterStatus.Status = "Status Tidak Terdefinisi"
		}

		data.WindStatus.Wind = wind
		switch {
		case wind < 6:
			data.WindStatus.Status = "Aman"
		case wind > 6 && water <= 15:
			data.WindStatus.Status = "Siaga"
		case wind > 15:
			data.WindStatus.Status = "Bahaya"
		default:
			data.WaterStatus.Status = "Status Tidak Terdefinisi"
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			log.Fatal("[error] occured while marshaling status data: ", err.Error())
		}

		if err = ioutil.WriteFile("data.json", jsonData, 0644); err != nil {
			log.Fatal("[error] occured while writing json data: ", err.Error())
		}

		time.Sleep(time.Second * 15)
	}
}

func AutoReloadWeb(w http.ResponseWriter, r *http.Request) {
	fileData, err := ioutil.ReadFile("data.json")

	if err != nil {
		log.Fatal("[error] error occured while reading data.json file: ", err.Error())
	}

	var statusData StatusData

	err = json.Unmarshal(fileData, &statusData)

	if err != nil {
		log.Fatal("[error] error occured while unmarshaling data.json file: ", err.Error())
	}

	waterValue := statusData.WaterStatus.Water
	waterStatus := statusData.WaterStatus.Status

	windValue := statusData.WindStatus.Wind
	windStatus := statusData.WindStatus.Status

	data := map[string]interface{}{
		"waterValue":  waterValue,
		"waterStatus": waterStatus,
		"windValue":   windValue,
		"windStatus":  windStatus,
	}

	tpl, err := template.ParseFiles("index.html")

	if err != nil {
		log.Fatal("[error] error occured while parsing html: ", err.Error())
	}

	tpl.Execute(w, data)
}
