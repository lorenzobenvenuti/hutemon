package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/lorenzobenvenuti/hutemon/weather"
)

func main() {
	flag.Parse()
	wuKey := os.Args[1]
	owmKey := os.Args[2]
	wp := make([]weather.WeatherProvider, 2)
	wp[0] = weather.NewWUndergroundWeatherProvider(wuKey)
	wp[1] = weather.NewOpenWeatherMapWeatherProvider(owmKey)
	w, err := weather.NewWeatherProviderChain(wp).GetWeather(os.Args[3])
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Printf("%#v\n", w)
}
