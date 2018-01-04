package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	key := os.Args[1]
	w, err := NewWUndergroundWeatherProvider(key).GetWeather(os.Args[2])
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Printf("%#v\n", w)
}
