package handlers

import (
	"testing"

	"github.com/lorenzobenvenuti/hutemon/sensor"
	"github.com/lorenzobenvenuti/hutemon/weather"
	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	wam := &WeatherAndMeasurement{}
	wam.Measurement = &sensor.Measurement{Temperature: 19.21, Humidity: 44.49}
	wam.Weather = &weather.Weather{Temperature: 11.32, Humidity: 87.87, Weather: "Sun"}
	assert.Equal(t, "44.5;19.2;87.9;11.3;Sun", getValue(wam), "Value sent is correctly formatted")
}

func TestNewIfttHandler(t *testing.T) {
	h := NewIftttHandler("key", "event")
	ih := h.(*iftttHandler)
	assert.Equal(t, "key", ih.apiKey, "IFTTT api key is correct")
	assert.Equal(t, "event", ih.eventName, "IFTTT event name is correct")
}
