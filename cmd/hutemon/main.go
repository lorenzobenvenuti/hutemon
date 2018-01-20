package main

import (
	"fmt"
	"time"

	"github.com/lorenzobenvenuti/hutemon/handlers"
	"github.com/lorenzobenvenuti/hutemon/weather"
)

func main() {
	// flag.Parse()
	// wuKey := os.Args[1]
	// owmKey := os.Args[2]
	// wp := make([]weather.Provider, 2)
	// wp[0] = weather.NewWUndergroundProvider(wuKey)
	// wp[1] = weather.NewOpenWeatherMapProvider(owmKey)
	// w, err := weather.NewProviderChain(wp).GetWeather(os.Args[3])
	// if err != nil {
	// 	glog.Fatal(err)
	// }
	// fmt.Printf("%#v\n", w)
	w := &weather.Weather{}
	h := make(map[string]handlers.Handler)
	h["a"] = handlers.NewDummyHandler("A", 1*time.Second)
	h["b"] = handlers.NewDummyHandler("B", 2*time.Second)
	h["c"] = handlers.NewDummyHandler("C", 3*time.Second)
	err := handlers.NewHandlerChain(h).Handle(w)
	fmt.Println(err)
}
