package weather

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyProvider struct {
	weather  *Weather
	err      error
	invoked  bool
	location string
}

func (dp *dummyProvider) GetWeather(location string) (*Weather, error) {
	dp.invoked = true
	dp.location = location
	return dp.weather, dp.err
}

func getWeather(w string, t float32, h float32) *Weather {
	return &Weather{Weather: w, Temperature: t, Humidity: h}
}

func TestProviderChainWhenFirstProviderReturnsWeather(t *testing.T) {
	p1 := &dummyProvider{weather: getWeather("Sun", 20, 50), err: nil}
	p2 := &dummyProvider{weather: getWeather("Clouds", 18, 48), err: nil}
	pc := &providerChain{providers: []Provider{p1, p2}}
	w, e := pc.GetWeather("Pisa")
	assert.True(t, p1.invoked, "First provider has been invoked")
	assert.Equal(t, "Pisa", p1.location, "First provider is invoked with right argument")
	assert.Nil(t, e, "Error is nil")
	assert.Equal(t, "Sun", w.Weather, "Weather is 'Sun'")
	assert.Equal(t, float32(20), w.Temperature, "Temperature is 20")
	assert.Equal(t, float32(50), w.Humidity, "humidity is 50")
	assert.False(t, p2.invoked, "Second provider is not invoked")
}

func TestProviderChainWhenFirstProviderReturnsError(t *testing.T) {
	p1 := &dummyProvider{weather: nil, err: errors.New("Error")}
	p2 := &dummyProvider{weather: getWeather("Clouds", 18, 48), err: nil}
	pc := &providerChain{providers: []Provider{p1, p2}}
	w, e := pc.GetWeather("Pisa")
	assert.True(t, p1.invoked, "First provider has been invoked")
	assert.Equal(t, "Pisa", p1.location, "First provider is invoked with right argument")
	assert.True(t, p2.invoked, "Second provider has been invoked")
	assert.Equal(t, "Pisa", p2.location, "Second provider is invoked with right argument")
	assert.Nil(t, e, "Error is nil")
	assert.Equal(t, "Clouds", w.Weather, "Weather is 'Clouds'")
	assert.Equal(t, float32(18), w.Temperature, "Temperature is 18")
	assert.Equal(t, float32(48), w.Humidity, "humidity is 48")
}

func TestProviderChainWhenBothProvidersReturnsError(t *testing.T) {
	p1 := &dummyProvider{weather: nil, err: errors.New("Error")}
	p2 := &dummyProvider{weather: nil, err: errors.New("Error")}
	pc := &providerChain{providers: []Provider{p1, p2}}
	w, e := pc.GetWeather("Pisa")
	assert.True(t, p1.invoked, "First provider has been invoked")
	assert.Equal(t, "Pisa", p1.location, "First provider is invoked with right argument")
	assert.True(t, p2.invoked, "First provider has been invoked")
	assert.Equal(t, "Pisa", p2.location, "First provider is invoked with right argument")
	assert.Nil(t, w, "Weather is nil")
	assert.Equal(t, "All providers failed to retrieve weather", e.Error(), "Error is not nil")
}
