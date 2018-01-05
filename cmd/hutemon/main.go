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
	key := os.Args[1]
	w, err := weather.NewWUndergroundWeatherProvider(key).GetWeather(os.Args[2])
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Printf("%#v\n", w)
}
