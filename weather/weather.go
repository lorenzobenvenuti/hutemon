package weather

import (
	"errors"

	"github.com/golang/glog"
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
			glog.Infof("Provider %v returned result: %v", wp, w)
			return w, nil
		}
		glog.Errorf("Provider %v returned an error: %s", wp, err)
	}
	return nil, errors.New("No provider did return weather")
}

func NewProviderChain(providers []Provider) Provider {
	return &providerChain{providers: providers}
}
