package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/subosito/gotenv"
)

type Geocode struct {
	Devicetime string  `json:"devicetime"`
	Fixtime    string  `json:"fixtime"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Altitude   float64 `json:"altitude"`
	Speed      float64 `json:"speed"`
}
type GeocodeReverse struct {
	Items []struct {
		Title           string `json:"title"`
		ID              string `json:"id"`
		Resulttype      string `json:"resultType"`
		Housenumbertype string `json:"houseNumberType"`
		Address         struct {
			Label       string `json:"label"`
			Countrycode string `json:"countryCode"`
			Countryname string `json:"countryName"`
			County      string `json:"county"`
			City        string `json:"city"`
			District    string `json:"district"`
			Street      string `json:"street"`
			Postalcode  string `json:"postalCode"`
			Housenumber string `json:"houseNumber"`
		} `json:"address"`
		Position struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"position"`
		Access []struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"access"`
		Distance int `json:"distance"`
		Mapview  struct {
			West  float64 `json:"west"`
			South float64 `json:"south"`
			East  float64 `json:"east"`
			North float64 `json:"north"`
		} `json:"mapView"`
	} `json:"items"`
}

type Trip struct {
	start Geocode
	end   Geocode
}

func init() {
	gotenv.Load()
}

func main() {
	file, _ := ioutil.ReadFile("geodata_quiz.json")

	g := []Geocode{}

	if err := json.Unmarshal([]byte(file), &g); err != nil {
		log.Printf("cannot convert: %+v", err)
	}

	g = reverseArray(g)

	log.Println(g[0])
	apikey := os.Getenv("api_key")
	tripMap := make(map[int]Trip)

	waintForStartTrip := true
	countTrip := 0
	var trip Trip
	for _, val := range g {

		if val.Speed > 2.00 && waintForStartTrip {
			waintForStartTrip = false
			trip.start = val
			trip.end = Geocode{}
			tripMap[countTrip] = trip
		}

		if val.Speed < 2.00 && !waintForStartTrip {
			waintForStartTrip = true
			trip.end = val
			tripMap[countTrip] = trip
			countTrip = countTrip + 1

			trip.start = Geocode{}
			trip.end = Geocode{}
		}
	}
	log.Println("The dataset have ", len(tripMap), " trip")

	for key, element := range tripMap {
		startTrip := getGeo(apikey, element.start)
		endTrip := getGeo(apikey, element.end)
		log.Println("********************* start of trip ", key, "***********************")
		log.Println("Trip:", key, "=>", "Strat at : ", startTrip.Items[0].Address.Housenumber, startTrip.Items[0].Address.Street, startTrip.Items[0].Address.District, startTrip.Items[0].Address.City, startTrip.Items[0].Address.Countryname)
		log.Println("Trip:", key, "=>", "End at : ", endTrip.Items[0].Address.Housenumber, endTrip.Items[0].Address.Street, endTrip.Items[0].Address.District, endTrip.Items[0].Address.City, endTrip.Items[0].Address.Countryname)
		log.Println("********************* end of trip ", key, "***********************")
	}
}

func getGeo(apikey string, data Geocode) GeocodeReverse {
	url := "https://revgeocode.search.hereapi.com/v1/revgeocode?apiKey=" + apikey + "&at=" + fmt.Sprint(data.Latitude) + "," + fmt.Sprint(data.Longitude)

	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	gecRev := GeocodeReverse{}
	if err := json.Unmarshal(body, &gecRev); err != nil {
		log.Println("connot Unmarshal", err)
	}
	return gecRev
}

func reverseArray(arr []Geocode) []Geocode {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
