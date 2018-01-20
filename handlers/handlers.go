package handlers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lorenzobenvenuti/hutemon/weather"
)

type Handler interface {
	Handle(w *weather.Weather) error
}

type handlerError struct {
	name string
	err  error
}

type multiError struct {
	errors map[string]error
}

func (me *multiError) Error() string {
	msgs := make([]string, len(me.errors))
	for k, e := range me.errors {
		msgs = append(msgs, fmt.Sprintf("%s -> %s", k, e))
	}
	return strings.Join(msgs, "\n")
}

type handlerChain struct {
	handers map[string]Handler
}

func (hc *handlerChain) Handle(w *weather.Weather) error {
	var wg sync.WaitGroup
	errors := make(map[string]error)
	queue := make(chan handlerError)
	wg.Add(len(hc.handers))
	go func() {
		for he := range queue {
			if he.err != nil {
				errors[he.name] = he.err
			}
			wg.Done()
		}
	}()
	for k, h := range hc.handers {
		go func(name string, handler Handler) {
			queue <- handlerError{name: name, err: handler.Handle(w)}
		}(k, h)
	}
	wg.Wait()
	if len(errors) == 0 {
		return nil
	}
	return &multiError{errors: errors}
}

func NewHandlerChain(handlers map[string]Handler) Handler {
	return &handlerChain{handers: handlers}
}
