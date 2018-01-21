package weather

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lorenzobenvenuti/hutemon/http"
	"github.com/sirupsen/logrus"
)

type wUndergroundProvider struct {
	apiKey           string
	httpClient       http.HttpClient
	jsonUnmarshaller http.JsonUnmarshaller
}

type wUndergroundResponse struct {
	CurrentObservation struct {
		Weather          string  `json:"weather"`
		TempC            float32 `json:"temp_c"`
		RelativeHumidity string  `json:"relative_humidity"`
	} `json:"current_observation"`
}

func (wu *wUndergroundProvider) GetWeather(location string) (*Weather, error) {
	url := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/IT/%s.json", wu.apiKey, location)
	wur := &wUndergroundResponse{}
	logrus.WithFields(logrus.Fields{
		"url":      url,
		"api_key":  wu.apiKey,
		"location": location,
	}).Debug("Performing request to Weather Underground")
	bytes, err := wu.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	err = wu.jsonUnmarshaller.Unmarshal(bytes, wur)
	if err != nil {
		return nil, err
	}
	h, err := strconv.ParseFloat(strings.Trim(wur.CurrentObservation.RelativeHumidity, "%"), 32)
	if err != nil {
		// TODO: loggare
		return nil, errors.New(fmt.Sprintf("Error parsing humidity (\"%s\")", wur.CurrentObservation.RelativeHumidity))
	}
	return &Weather{Weather: wur.CurrentObservation.Weather, Temperature: wur.CurrentObservation.TempC, Humidity: float32(h)}, nil
}

func NewWUndergroundProvider(apiKey string) Provider {
	return &wUndergroundProvider{
		apiKey:           apiKey,
		httpClient:       http.NewHttpClient(10 * time.Second),
		jsonUnmarshaller: http.NewJsonUnmarshaller(),
	}
}
