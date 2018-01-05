package sensor

import dht "github.com/d2r2/go-dht"

type dht22Sensor struct{}

func (d *dht22Sensor) Measure() (*Measurement, error) {
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, 4, true, 10)
	if err != nil {
		return nil, err
	}
	return &Measurement{Temperature: temperature, Humidity: humidity}, nil
}

func NewDHT22Sensor() Sensor {
	return &dht22Sensor{}
}
