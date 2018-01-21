package weather

import (
	"fmt"
	"time"

	"github.com/lorenzobenvenuti/hutemon/http"
)

type openWeatherMapProvider struct {
	apiKey           string
	httpClient       http.HttpClient
	jsonUnmarshaller http.JsonUnmarshaller
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
	bytes, err := wp.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	err = wp.jsonUnmarshaller.Unmarshal(bytes, owr)
	if err != nil {
		return nil, err
	}
	return &Weather{Weather: owr.Weather[0].Main, Temperature: owr.Main.Temp, Humidity: float32(owr.Main.Humidity)}, nil
}

func NewOpenWeatherMapProvider(apiKey string) Provider {
	return &openWeatherMapProvider{
		apiKey:           apiKey,
		httpClient:       http.NewHttpClient(10 * time.Second),
		jsonUnmarshaller: http.NewJsonUnmarshaller(),
	}
}
