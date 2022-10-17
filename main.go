package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Data struct {
	Status `json:"status"`
}

type StatusBencana struct {
	StatusWater string
	StatusWind  string
}

func updateData() {
	for {
		var data = Data{
			Status: Status{},
		}

		minValue := 1
		maxValue := 100

		data.Status.Water = rand.Intn(maxValue-minValue) + minValue
		data.Status.Wind = rand.Intn(maxValue-minValue) + minValue

		b, err := json.MarshalIndent(&data, "", " ")

		if err != nil {
			fmt.Println("Error marshal indent:", err)
			return
		}

		err = os.WriteFile("file.json", b, 0644)

		if err != nil {
			fmt.Println("Error writing file:", err)
		}

		fmt.Println("Menunggu 15 detik")
		time.Sleep(15 * time.Second)
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())

	go updateData()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var status = StatusBencana{}

		templ, err := template.ParseFiles("index.html")

		if err != nil {
			fmt.Println("Error parsing file:", err)
			return
		}

		var data = Data{Status: Status{}}

		b, err := os.ReadFile("file.json")

		if err != nil {
			fmt.Println("Error read file:", err)
			return
		}

		err = json.Unmarshal(b, &data)

		if err != nil {
			fmt.Println("Error unmarshal:", err)
			return
		}

		if data.Status.Water < 5 {
			status.StatusWater = "aman"
		} else if data.Status.Water > 5 && data.Status.Water < 9 {
			status.StatusWater = "siaga"
		} else if data.Status.Water > 8 {
			status.StatusWater = "bahaya"
		}

		if data.Status.Wind < 6 {
			status.StatusWind = "aman"
		} else if data.Status.Wind > 6 && data.Status.Wind < 16 {
			status.StatusWind = "siaga"
		} else if data.Status.Wind > 15 {
			status.StatusWind = "bahaya"
		}

		value := map[string]any{
			"waterValue":  data.Status.Water,
			"waterStatus": status.StatusWater,
			"windValue":   data.Status.Wind,
			"windStatus":  status.StatusWind,
		}

		err = templ.ExecuteTemplate(w, "index.html", value)

		if err != nil {
			fmt.Println("Error executing template:", err)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
