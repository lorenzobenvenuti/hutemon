package handlers

import (
	"fmt"
	"testing"
	"time"

	"github.com/lorenzobenvenuti/hutemon/weather"
	"github.com/stretchr/testify/assert"
)

type dummyHandler struct {
	name    string
	timeout time.Duration
	failing bool
}

func (dh *dummyHandler) Handle(w *weather.Weather) error {
	time.Sleep(dh.timeout)
	if dh.failing {
		return fmt.Errorf("Error %s", dh.name)
	}
	return nil
}

func newDummyHandler(name string, timeout time.Duration, failing bool) Handler {
	return &dummyHandler{name: name, timeout: timeout, failing: failing}
}

func TestHandlerChainWithErrors(t *testing.T) {
	w := &weather.Weather{}
	h := make(map[string]Handler)
	h["a"] = newDummyHandler("A", 10*time.Millisecond, true)
	h["b"] = newDummyHandler("B", 20*time.Millisecond, true)
	h["c"] = newDummyHandler("C", 30*time.Millisecond, false)
	h["d"] = newDummyHandler("D", 40*time.Millisecond, true)
	err := NewHandlerChain(h).Handle(w)
	me, ok := err.(*multiError)
	assert.True(t, ok, "Error must be a multiError")
	assert.Equal(t, 3, len(me.errors), "Three errors expected")
	assert.Contains(t, me.errors, "a")
	assert.Equal(t, me.errors["a"].Error(), "Error A")
	assert.Contains(t, me.errors, "b")
	assert.Equal(t, me.errors["b"].Error(), "Error B")
	assert.Contains(t, me.errors, "a")
	assert.Equal(t, me.errors["d"].Error(), "Error D")
}

func TestHandlerChainWithoutErrors(t *testing.T) {
	w := &weather.Weather{}
	h := make(map[string]Handler)
	h["a"] = newDummyHandler("A", 10*time.Millisecond, false)
	h["b"] = newDummyHandler("B", 20*time.Millisecond, false)
	h["c"] = newDummyHandler("C", 30*time.Millisecond, false)
	err := NewHandlerChain(h).Handle(w)
	assert.Nil(t, err, "No error should be returned")
}
