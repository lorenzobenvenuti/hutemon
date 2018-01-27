package handlers

import "github.com/lorenzobenvenuti/ifttt"

type iftttHandler struct {
	apiKey    string
	eventName string
}

func (ih *iftttHandler) Handle(wam *WeatherAndMeasurement) error {
	return ifttt.NewIftttClient(ih.apiKey).Trigger(ih.eventName, []string{})
}
