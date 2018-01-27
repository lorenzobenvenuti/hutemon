package handlers

import "github.com/sirupsen/logrus"

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
