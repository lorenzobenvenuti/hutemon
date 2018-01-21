package handlers

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dummyHandler struct {
	name    string
	timeout time.Duration
	failing bool
}

func (dh *dummyHandler) Handle(w *WeatherAndMeasurement) error {
	time.Sleep(dh.timeout)
	if dh.failing {
		return fmt.Errorf("Error %s", dh.name)
	}
	return nil
}

func newDummyHandler(name string, timeout time.Duration, failing bool) Handler {
	return &dummyHandler{name: name, timeout: timeout, failing: failing}
}

func TestMultiError(t *testing.T) {
	errs := make(map[string]error)
	errs["a"] = errors.New("Error A")
	errs["b"] = errors.New("Error B")
	errs["c"] = errors.New("Error C")
	me := &multiError{errors: errs}
	str := me.Error()
	lines := strings.Split(str, "\n")
	assert.Equal(t, 3, len(lines), "Three lines expected")
	assert.Contains(t, lines, "a -> Error A")
	assert.Contains(t, lines, "b -> Error B")
	assert.Contains(t, lines, "c -> Error C")
}

func TestHandlerChainWithErrors(t *testing.T) {
	w := &WeatherAndMeasurement{}
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
	assert.Equal(t, "Error A", me.errors["a"].Error())
	assert.Contains(t, me.errors, "b")
	assert.Equal(t, "Error B", me.errors["b"].Error())
	assert.Contains(t, me.errors, "a")
	assert.Equal(t, "Error D", me.errors["d"].Error())
}

func TestHandlerChainWithoutErrors(t *testing.T) {
	w := &WeatherAndMeasurement{}
	h := make(map[string]Handler)
	h["a"] = newDummyHandler("A", 10*time.Millisecond, false)
	h["b"] = newDummyHandler("B", 20*time.Millisecond, false)
	h["c"] = newDummyHandler("C", 30*time.Millisecond, false)
	err := NewHandlerChain(h).Handle(w)
	assert.Nil(t, err, "No error should be returned")
}
