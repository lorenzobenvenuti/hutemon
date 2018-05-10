# hutemon - HUmidity and TEmperature MONitor

A tool to monitor humidity and temperature of my basement, using a Raspberry Pi and a DHT22 sensor.

# TODO

* Rename `handlerChain` (is not a chain, handlers are run in parallel)
* Implement InfluxDb client
* Use goroutines to retrieve measurement and weather in parallel
* Use `mock` (from `testify`) instead of stubs
* Parse command line arguments
* Test package `http`
* Move dummy sensor to tests of `main` package
* Logging
* Test factory methods
* `errors` package for error context
* Log "technical" errors and return user friendly messages
