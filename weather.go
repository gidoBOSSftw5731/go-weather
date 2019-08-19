package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	endpoint = "https://api.openweathermap.org/data/2.5/weather?%v=%v,%v&appid=%v&units=metric"
)

//Weather is a struct that contains json information of weather for a specified location
type Weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func keyRemover(err error, key string) error {
	stringErr := strings.Replace(fmt.Sprint(err), key, "${KEY}", -1)

	return fmt.Errorf(stringErr)
}

//CurrentWeather returns a struct with weather data.
func CurrentWeather(location, key string) (Weather, error) {
	var w Weather

	println(location)

	isZip, err := regexp.MatchString("\\d{5}(?:[-s]\\d{4})?", location)
	if err != nil {
		err = keyRemover(err, key)
		return w, err
	}

	if isZip {
		country := "us"
		zip := "94040"

		forgein, err := regexp.MatchString("/[A-Za-z]{2}+/", location)
		if err != nil {
			err = keyRemover(err, key)
			return w, err
		}

		if forgein {
			parts := strings.Split(location, " ")
			if len(parts) == 1 {
				return w, fmt.Errorf("No space between the country code and zip")
			}
			for i := 0; i < len(parts); i++ {
				isCountry, _ := regexp.MatchString("/^[A-Za-z]{2}+$/", parts[i])
				if isCountry {
					country = parts[i]
					break
				}
			}
		}

		parts := strings.Split(location, " ")
		for i := 0; i < len(parts); i++ {
			zipFound, _ := regexp.MatchString("\\d{5}(?:[-\\s]\\d{4})?", parts[i])
			if zipFound {
				zip = parts[i]
				break
			}
		}

		url := fmt.Sprintf(endpoint, "zip", zip, country, key)

		req, err := http.Get(url)
		if err != nil {
			err = keyRemover(err, key)
			return w, err
		}
		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			err = keyRemover(err, key)
			return w, err
		}

		err = json.Unmarshal(body, &w)

		//fmt.Println(string(body), "\n", url)
	} else {
		country := "us"
		var city string

		forgein, err := regexp.MatchString("/[A-Za-z]{2}+/", location)
		if err != nil {
			err = keyRemover(err, key)
			return w, err
		}

		if forgein {
			parts := strings.Split(location, " ")
			if len(parts) == 1 {
				return w, fmt.Errorf("No space between the country code and zip")
			}
			for i := 0; i < len(parts); i++ {
				isCountry, _ := regexp.MatchString("/^[A-Za-z]{2}+$/", parts[i])
				if isCountry {
					country = parts[i]
					break
				} else {
					city += parts[i]
				}
			}
		}

		url := fmt.Sprintf(endpoint, "q", city, country, key)

		req, err := http.Get(url)
		if err != nil {
			err = keyRemover(err, key)
			return w, err
		}
		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			err = keyRemover(err, key)
			return w, err
		}

		err = json.Unmarshal(body, &w)

		//fmt.Println(string(body), "\n", url)

	}

	return w, nil

}
