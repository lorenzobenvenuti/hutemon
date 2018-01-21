package sensor

import dht "github.com/d2r2/go-dht"

type dht22Sensor struct {
	pin   int
	retry int
}

func (d *dht22Sensor) Measure() (*Measurement, error) {
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, d.pin, true, d.retry)
	if err != nil {
		return nil, err
	}
	return &Measurement{Temperature: temperature, Humidity: humidity}, nil
}

func NewDHT22Sensor(pin int, retry int) Sensor {
	return &dht22Sensor{pin: pin, retry: retry}
}
