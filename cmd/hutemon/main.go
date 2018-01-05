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
	wp := make([]weather.Provider, 2)
	wp[0] = weather.NewWUndergroundProvider(wuKey)
	wp[1] = weather.NewOpenWeatherMapProvider(owmKey)
	w, err := weather.NewProviderChain(wp).GetWeather(os.Args[3])
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Printf("%#v\n", w)
}
