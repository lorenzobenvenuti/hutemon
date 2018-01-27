package handlers

import (
	"fmt"

	"github.com/lorenzobenvenuti/ifttt"
)

type iftttHandler struct {
	apiKey    string
	eventName string
}

func getValue(wam *WeatherAndMeasurement) string {
	return fmt.Sprintf(
		"%.1f;%.1f;%.1f;%.1f;%s",
		wam.Measurement.Humidity,
		wam.Measurement.Temperature,
		wam.Weather.Humidity,
		wam.Weather.Temperature,
		wam.Weather.Weather,
	)
}

func (ih *iftttHandler) Handle(wam *WeatherAndMeasurement) error {
	return ifttt.NewIftttClient(ih.apiKey).Trigger(ih.eventName, []string{getValue(wam)})
}

func NewIftttHandler(apiKey string, eventName string) Handler {
	return &iftttHandler{apiKey: apiKey, eventName: eventName}
}
