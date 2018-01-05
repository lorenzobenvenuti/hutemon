package sensor

type Measurement struct {
	Temperature float32
	Humidity    float32
}

type Sensor interface {
	Measure() (*Measurement, error)
}

type dummySensor struct {
	temperature float32
	humidity    float32
}

func (ds *dummySensor) Measure() (*Measurement, error) {
	return &Measurement{Temperature: ds.temperature, Humidity: ds.humidity}, nil
}

func NewDummySensor(t float32, h float32) Sensor {
	return &dummySensor{temperature: t, humidity: h}
}
