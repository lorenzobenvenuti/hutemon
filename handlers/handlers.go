package handlers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lorenzobenvenuti/hutemon/sensor"
	"github.com/lorenzobenvenuti/hutemon/weather"
	"github.com/sirupsen/logrus"
)

type WeatherAndMeasurement struct {
	Weather     *weather.Weather
	Measurement *sensor.Measurement
}

type Handler interface {
	Handle(wam *WeatherAndMeasurement) error
}

type handlerError struct {
	name string
	err  error
}

type multiError struct {
	errors map[string]error
}

func (me *multiError) Error() string {
	msgs := make([]string, 0)
	for k, e := range me.errors {
		msgs = append(msgs, fmt.Sprintf("%s -> %s", k, e))
	}
	return strings.Join(msgs, "\n")
}

type handlerChain struct {
	handers map[string]Handler
}

func (hc *handlerChain) Handle(wam *WeatherAndMeasurement) error {
	var wg sync.WaitGroup
	errors := make(map[string]error)
	queue := make(chan handlerError)
	wg.Add(len(hc.handers))
	go func() {
		for he := range queue {
			if he.err != nil {
				errors[he.name] = he.err
			}
			wg.Done()
		}
	}()
	for k, h := range hc.handers {
		go func(name string, handler Handler) {
			queue <- handlerError{name: name, err: handler.Handle(wam)}
		}(k, h)
	}
	wg.Wait()
	if len(errors) == 0 {
		return nil
	}
	return &multiError{errors: errors}
}

func NewHandlerChain(handlers map[string]Handler) Handler {
	return &handlerChain{handers: handlers}
}

type loggingHandler struct{}

func (lh *loggingHandler) Handle(wam *WeatherAndMeasurement) error {
	logrus.WithFields(logrus.Fields{
		"weather":     wam.Weather,
		"measurement": wam.Measurement,
	}).Info("Received weather and measurement")
	return nil
}

func NewLoggingHandler() Handler {
	return &loggingHandler{}
}
