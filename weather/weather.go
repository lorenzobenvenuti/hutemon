package weather

import (
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Weather struct {
	Weather     string
	Temperature float32
	Humidity    float32
}

type Provider interface {
	GetWeather(location string) (*Weather, error)
}

type providerChain struct {
	providers []Provider
}

func (wpc *providerChain) GetWeather(location string) (*Weather, error) {
	for _, wp := range wpc.providers {
		w, err := wp.GetWeather(location)
		if err == nil {
			logrus.WithFields(logrus.Fields{
				"provider": reflect.TypeOf(wp),
				"result":   w,
			}).Debug("Weather correctly retrieved")
			return w, nil
		}
		logrus.WithFields(logrus.Fields{
			"provider": wp,
			"error":    err,
		}).Warn("Error retrieving weather")
	}
	return nil, errors.New("All providers failed to retrieve weather")
}

func NewProviderChain(providers []Provider) Provider {
	return &providerChain{providers: providers}
}
