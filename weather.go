package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
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

	client := http.Client{Timeout: time.Second * 10}
	url := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/IT/%s.json", wu.apiKey, location)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	wur := &wUndergroundResponse{}
	err = json.Unmarshal(body, wur)
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

type weatherProviderChain struct {
	weatherProviders []WeatherProvider
}

func (wpc *weatherProviderChain) GetWeather(location string) (*Weather, error) {
	for _, wp := range wpc.weatherProviders {
		w, err := wp.GetWeather(location)
		if err != nil {
			glog.Errorf("Provider %v returned an error: %s", w, err)
			return w, nil
		}
	}
	return nil, errors.New("No provider did return weather")
}

func NewWeatherProviderChain(weatherProviders []WeatherProvider) WeatherProvider {
	return &weatherProviderChain{weatherProviders: weatherProviders}
}
