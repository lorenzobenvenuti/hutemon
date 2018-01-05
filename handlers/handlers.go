package handlers

import "github.com/lorenzobenvenuti/hutemon/weather"

type Handler interface {
	Handle(w weather.Weather) error
}
