package weather

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/lorenzobenvenuti/hutemon/http"
)

// Weather is a struct representing the weather
type Weather struct {
	Weather     string
	Temperature float32
	Humidity    float32
}

// WeatherProvider can load the weather
type WeatherProvider interface {
	GetWeather(location string) (*Weather, error)
}

type wUndergroundWeatherProvider struct {
	apiKey string
}

type wUndergroundResponse struct {
	CurrentObservation struct {
		Weather          string  `json:"weather"`
		TempC            float32 `json:"temp_c"`
		RelativeHumidity string  `json:"relative_humidity"`
	} `json:"current_observation"`
}

func (wu *wUndergroundWeatherProvider) GetWeather(location string) (*Weather, error) {
	url := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/IT/%s.json", wu.apiKey, location)
	wur := &wUndergroundResponse{}
	err := http.GetAndUnmarshal(url, wur)
	if err != nil {
		return nil, err
	}

	h, err := strconv.ParseFloat(strings.Trim(wur.CurrentObservation.RelativeHumidity, "%"), 32)
	if err != nil {
		return nil, err
	}
	return &Weather{Weather: wur.CurrentObservation.Weather, Temperature: wur.CurrentObservation.TempC, Humidity: float32(h)}, nil
}

func NewWUndergroundWeatherProvider(apiKey string) WeatherProvider {
	return &wUndergroundWeatherProvider{apiKey: apiKey}
}

type openWeatherMapWeatherProvider struct {
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

func (wp *openWeatherMapWeatherProvider) GetWeather(location string) (*Weather, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s&units=metric", wp.apiKey, location)
	owr := &openWeatherMapResponse{}
	err := http.GetAndUnmarshal(url, owr)
	if err != nil {
		return nil, err
	}
	return &Weather{Weather: owr.Weather[0].Main, Temperature: owr.Main.Temp, Humidity: float32(owr.Main.Humidity)}, nil
}

func NewOpenWeatherMapWeatherProvider(apiKey string) WeatherProvider {
	return &openWeatherMapWeatherProvider{apiKey: apiKey}
}

type weatherProviderChain struct {
	weatherProviders []WeatherProvider
}

func (wpc *weatherProviderChain) GetWeather(location string) (*Weather, error) {
	for _, wp := range wpc.weatherProviders {
		w, err := wp.GetWeather(location)
		if err == nil {
			glog.Infof("Provider %v returned result: %v", wp, w)
			return w, nil
		}
		glog.Errorf("Provider %v returned an error: %s", wp, err)
	}
	return nil, errors.New("No provider did return weather")
}

func NewWeatherProviderChain(weatherProviders []WeatherProvider) WeatherProvider {
	return &weatherProviderChain{weatherProviders: weatherProviders}
}
