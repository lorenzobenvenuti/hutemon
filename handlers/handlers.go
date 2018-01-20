package handlers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lorenzobenvenuti/hutemon/weather"
)

type Handler interface {
	Handle(w *weather.Weather) error
}

type dummyHandler struct {
	name    string
	timeout time.Duration
}

func (dh *dummyHandler) Handle(w *weather.Weather) error {
	time.Sleep(dh.timeout)
	return fmt.Errorf("Error %s", dh.name)
}

func NewDummyHandler(name string, timeout time.Duration) Handler {
	return &dummyHandler{name: name, timeout: timeout}
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
		msgs = append(msgs, fmt.Sprintf("%v -> %s", k, e))
	}
	return strings.Join(msgs, "\n")
}

type handlerChain struct {
	handers map[string]Handler
}

func (hc *handlerChain) Handle(w *weather.Weather) error {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	errors := make(map[string]error)
	queue := make(chan handlerError, 1)
	wg.Add(len(hc.handers))
	for k, h := range hc.handers {
		go func(name string, handler Handler) {
			queue <- handlerError{name: name, err: handler.Handle(w)}
		}(k, h)
		go func() {
			for he := range queue {
				if he.err != nil {
					mutex.Lock()
					errors[he.name] = he.err
					mutex.Unlock()
				}
				wg.Done()
			}
		}()
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
