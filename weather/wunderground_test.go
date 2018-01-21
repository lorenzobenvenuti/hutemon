package weather

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHttpClient struct {
	bytes []byte
	err   error
	url   string
}

func (hc *mockHttpClient) Get(url string) ([]byte, error) {
	hc.url = url
	return hc.bytes, hc.err
}

type mockJsonUnmarshaller struct {
	invoked     bool
	weather     string
	temperature float32
	humidity    float32
}

func (ju *mockJsonUnmarshaller) Unmarshal(bytes []byte, v interface{}) error {
	ju.invoked = true
	if len(bytes) == 2 && bytes[0] == 1 && bytes[1] == 2 {
		w := v.(*wUndergroundResponse)
		w.CurrentObservation.Weather = ju.weather
		w.CurrentObservation.TempC = ju.temperature
		w.CurrentObservation.RelativeHumidity = fmt.Sprintf("%f", ju.humidity)
		return nil
	}
	if len(bytes) == 2 && bytes[0] == 3 && bytes[1] == 4 {
		w := v.(*openWeatherMapResponse)
		w.Main.Temp = ju.temperature
		w.Main.Humidity = int(ju.humidity)
		//w.Weather = make([]openWeatherMapResponse.Main)
		return nil
	}
	return errors.New("Error during unmarshalling")
}

func TestWUndergroundProviderOk(t *testing.T) {
	hc := &mockHttpClient{bytes: []byte{1, 2}}
	wUndergroundProvider := &wUndergroundProvider{
		apiKey:           "my-api-key",
		httpClient:       hc,
		jsonUnmarshaller: &mockJsonUnmarshaller{weather: "Sun", temperature: 20, humidity: 50},
	}
	resp, err := wUndergroundProvider.GetWeather("Pisa")
	assert.Equal(t, "Sun", resp.Weather, "Weather is 'Sun'")
	assert.Equal(t, float32(20), resp.Temperature, "Temperature is 20")
	assert.Equal(t, float32(50), resp.Humidity, "Humidity is 50")
	assert.Equal(t, "http://api.wunderground.com/api/my-api-key/conditions/q/IT/Pisa.json", hc.url, "Url string contains api key and location")
	assert.Nil(t, err, "Error is nil")
}

func TestWUndergroundProviderWhenHttpClientReturnsError(t *testing.T) {
	hc := &mockHttpClient{err: errors.New("HTTP Error")}
	ju := &mockJsonUnmarshaller{weather: "Sun", temperature: 20, humidity: 50}
	wUndergroundProvider := &wUndergroundProvider{
		apiKey:           "my-api-key",
		httpClient:       hc,
		jsonUnmarshaller: ju,
	}
	resp, err := wUndergroundProvider.GetWeather("Pisa")
	assert.Nil(t, resp, "Response is nil")
	assert.Equal(t, "HTTP Error", err.Error(), "Error is not null")
	assert.False(t, ju.invoked, "Marshaller is not invoked")
}

func TestWUndergroundProviderWhenUnmsarhallerReturnsError(t *testing.T) {
	hc := &mockHttpClient{bytes: []byte{3}}
	wUndergroundProvider := &wUndergroundProvider{
		apiKey:           "my-api-key",
		httpClient:       hc,
		jsonUnmarshaller: &mockJsonUnmarshaller{weather: "Sun", temperature: 20, humidity: 50},
	}
	resp, err := wUndergroundProvider.GetWeather("Pisa")
	assert.Nil(t, resp, "Response is nil")
	assert.Equal(t, "Error during unmarshalling", err.Error(), "Error is not null")
	assert.Equal(t, "http://api.wunderground.com/api/my-api-key/conditions/q/IT/Pisa.json", hc.url, "Url string contains api key and location")
}
