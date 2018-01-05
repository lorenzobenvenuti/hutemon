package weather

import (
	"fmt"

	"github.com/lorenzobenvenuti/hutemon/http"
)

type openWeatherMapProvider struct {
	apiKey string
}

type openWeatherMapResponse struct {
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp     float32 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
}

func (wp *openWeatherMapProvider) GetWeather(location string) (*Weather, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s&units=metric", wp.apiKey, location)
	owr := &openWeatherMapResponse{}
	err := http.GetAndUnmarshal(url, owr)
	if err != nil {
		return nil, err
	}
	return &Weather{Weather: owr.Weather[0].Main, Temperature: owr.Main.Temp, Humidity: float32(owr.Main.Humidity)}, nil
}

func NewOpenWeatherMapProvider(apiKey string) Provider {
	return &openWeatherMapProvider{apiKey: apiKey}
}
