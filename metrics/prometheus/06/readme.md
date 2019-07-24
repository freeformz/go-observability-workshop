# Middleware

The package `promhttp` provides various standard lib compatible middleware that returns an instrumented [`http.Handler`](https://golang.org/pkg/net/http/#Handler). The middleware you choose depends on the type of metrics you are collecting:

* [`InstrumentHandlerCounter`](https://godoc.org/github.com/prometheus/client_golang/prometheus/promhttp#InstrumentHandlerCounter) for incrementing a counter.
* [`InstrumentHandlerDuration`](https://godoc.org/github.com/prometheus/client_golang/prometheus/promhttp#InstrumentHandlerDuration) for updating an [`ObserverVec`](https://godoc.org/github.com/prometheus/client_golang/prometheus#ObserverVec) with request duration information.
* [`InstrumentHandlerInFlight`](https://godoc.org/github.com/prometheus/client_golang/prometheus/promhttp#InstrumentHandlerInFlight) for setting a gauge to the number of inflight requests being handled.
* [`InstrumentHandlerRequestSize`](https://godoc.org/github.com/prometheus/client_golang/prometheus/promhttp#InstrumentHandlerRequestSize) for updating an [`ObserverVec`](https://godoc.org/github.com/prometheus/client_golang/prometheus#ObserverVec) with request size information.
* [`InstrumentHandlerTimeToWriteHeader`](https://godoc.org/github.com/prometheus/client_golang/prometheus/promhttp#InstrumentHandlerTimeToWriteHeader) for updating an [`ObserverVec`](https://godoc.org/github.com/prometheus/client_golang/prometheus#ObserverVec) with request time to write the first header.

Like all middleware they can be used in a chain:

```go
InstrumentHandlerInFlight(inFlightGauge,
  InstrumentHandlerDuration(durationVec,
    InstrumentHandlerRequestSize(sizeVec,
      myHandler,
    ),
  ),
)
```

## Exercise

Starting with the code from the last exercise, rewrite the server to the `InstrumentHandlerDuration` handler from the `promhttp` package.

Make sure to shutdown the prometheus server and remove the `data` directory to nuk any old data points. Follow the instructions from the last exercise to restart prometheus and the 2 instances of hey.

Run the queries from the last exercise and see what, if anything has changed.

Bonus activity: Utilize multiple handlers, which will require additional metrics.
