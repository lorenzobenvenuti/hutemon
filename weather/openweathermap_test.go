package weather

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type openWeatherMapJsonUnmarshaller struct {
	invoked     bool
	weather     string
	temperature float32
	humidity    int
}

func (ju *openWeatherMapJsonUnmarshaller) Unmarshal(bytes []byte, v interface{}) error {
	ju.invoked = true
	if len(bytes) == 2 && bytes[0] == 1 && bytes[1] == 2 {
		w := v.(*openWeatherMapResponse)
		w.Main.Humidity = ju.humidity
		w.Main.Temp = ju.temperature
		// json annotation required?! FIXME!?
		w.Weather = []struct {
			Main string `json:"main"`
		}{
			{Main: "Sun"},
		}
		return nil
	}
	return errors.New("Error during unmarshalling")
}

func TestOpenWeatherMapProviderOk(t *testing.T) {
	hc := &mockHttpClient{bytes: []byte{1, 2}}
	owp := &openWeatherMapProvider{
		apiKey:           "my-api-key",
		httpClient:       hc,
		jsonUnmarshaller: &openWeatherMapJsonUnmarshaller{weather: "Sun", temperature: 20, humidity: 50},
	}
	resp, err := owp.GetWeather("Pisa")
	assert.Equal(t, "Sun", resp.Weather, "Weather is 'Sun'")
	assert.Equal(t, float32(20), resp.Temperature, "Temperature is 20")
	assert.Equal(t, float32(50), resp.Humidity, "Humidity is 50")
	assert.Equal(t, "http://api.openweathermap.org/data/2.5/weather?appid=my-api-key&q=Pisa&units=metric", hc.url, "Url string contains api key and location")
	assert.Nil(t, err, "Error is nil")
}

func TestOpenWeatherMapProviderWhenHttpClientReturnsError(t *testing.T) {
	hc := &mockHttpClient{err: errors.New("HTTP Error")}
	ju := &openWeatherMapJsonUnmarshaller{weather: "Sun", temperature: 20, humidity: 50}
	owp := &openWeatherMapProvider{
		apiKey:           "my-api-key",
		httpClient:       hc,
		jsonUnmarshaller: ju,
	}
	resp, err := owp.GetWeather("Pisa")
	assert.Nil(t, resp, "Response is nil")
	assert.Equal(t, "HTTP Error", err.Error(), "Error is not null")
	assert.False(t, ju.invoked, "Marshaller is not invoked")
}

func TestOpenWeatherMapProviderWhenUnmarshallerReturnsError(t *testing.T) {
	hc := &mockHttpClient{bytes: []byte{3}}
	owp := &openWeatherMapProvider{
		apiKey:           "my-api-key",
		httpClient:       hc,
		jsonUnmarshaller: &openWeatherMapJsonUnmarshaller{weather: "Sun", temperature: 20, humidity: 50},
	}
	resp, err := owp.GetWeather("Pisa")
	assert.Nil(t, resp, "Response is nil")
	assert.Equal(t, "Error during unmarshalling", err.Error(), "Error is not null")
}

func TestNewOpenWeatherMapProvider(t *testing.T) {
	p := NewOpenWeatherMapProvider("key")
	owp := p.(*openWeatherMapProvider)
	assert.Equal(t, "key", owp.apiKey, "Api key must be the one passed to factory method")
}
