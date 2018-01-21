package main

import (
	"os"

	"github.com/lorenzobenvenuti/hutemon/handlers"
	"github.com/lorenzobenvenuti/hutemon/sensor"
	"github.com/lorenzobenvenuti/hutemon/weather"
	"github.com/sirupsen/logrus"
)

func getProvider(args []string) weather.Provider {
	wuKey := os.Args[1]
	owmKey := os.Args[2]
	wp := make([]weather.Provider, 2)
	wp[0] = weather.NewWUndergroundProvider(wuKey)
	wp[1] = weather.NewOpenWeatherMapProvider(owmKey)
	return weather.NewProviderChain(wp)
}

func getHandler(args []string) handlers.Handler {
	h := make(map[string]handlers.Handler)
	h["logging"] = handlers.NewLoggingHandler()
	return handlers.NewHandlerChain(h)
}

func main() {
	m, err := sensor.NewDHT22Sensor(4, 10).Measure()
	if err != nil {
		logrus.Fatal(err)
	}
	w, err := getProvider(os.Args).GetWeather(os.Args[3])
	if err != nil {
		logrus.Error("Cannot retrieve weather")
	}
	wam := &handlers.WeatherAndMeasurement{Weather: w, Measurement: m}
	getHandler(os.Args).Handle(wam)
}
