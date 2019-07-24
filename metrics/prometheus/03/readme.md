# Vectored Prometheus Application Metrics

In the last section we re-instrumented our app with the Prometheus `Counter` type.

During the [intro](../intro.md) we learned how different metrics can contain multiple time series. In the Go client these are implemented via "Vectors".

Vectors are the application of labels (and values) to a named metric.

```go
failures := prometheus.NewCounterVec(
  prometheus.CounterOpts{
    Name: "failures_total",
    Help: "Number of failures.",
  },
  []string{"host"},
)
```

The above example creates a Counter named `failures_total` that has one label key `host`.

Each of the different types have concrete implementations named `<T>Vec`:

* [CounterVec](https://godoc.org/github.com/prometheus/client_golang/prometheus#CounterVec)
* [GuageVec](https://godoc.org/github.com/prometheus/client_golang/prometheus#GaugeVec)
* [HistogramVec](https://godoc.org/github.com/prometheus/client_golang/prometheus#HistogramVec)
* [SummaryVec](https://godoc.org/github.com/prometheus/client_golang/prometheus#SummaryVec)

`Vec` style values created with a set of labels, require the label's values to be specified when recording a value.

`Vec` style types provide APIs for filling in the label values. These APIs either return errors (`GetMetricWith` & `GetMetricWithLabelValues`)
 or panic (`With` & `WithLabelValues`).

Additionally `Vec` style types have the `CurryWith` (returns error) and `MustCurryWith` (panic) that create `Vec` style values with partially filled label data.

Example:

```go
failures := prometheus.NewCounterVec(
  prometheus.CounterOpts{
    Name: "failures_total",
    Help: "Number of failures.",
  },
  []string{"host","type"},
)
// failures{}

host := os.Getenv("HOSTNAME")
if host == "" {
  host = "unknown"
}
failures = failures.MustCurryWith(prometheus.Labels{"host":host})
// failures{host="www.google.com"}

...

if err != nil {
  // would panic if not the right number of label values
  failures.WithLabelValues("some_error").Add(1)
}
// failures{host="www.google.com", type="some_error"} == 1
```

Note: You can use `WithLabelValues` and `GetMetricWithLabelValues` APIs to create metrics without any samples or observations.
This is a common pattern to avoid the pitfall of missing time series in metrics.

## Exercise

In the last exercise we ended up with metrics without labels, nor was the expvar program info converted.

We used two different metrics `http_requests_total` and `http_errors_total`.

Remember the advice from the [intro](../intro.md) though:

```text
Metric names specify the general aspect of a system that is being measured.
```

So what we want to do instead is to model these metrics via metric streams using labels.

Convert the two prometheus counters into a single, vectored counter with a `code` label.

When done, use curl to view '/metrics' and `hey` to induce load.

Bonus activity if you are done quickly:
Convert the expvar `Port` information into a `program_info` `Gauge` with label/values that expose the port.

## Give it a try

```console
$ go run server.go &
$ hey -c 1 -z 60m http://localhost:8080/ &
$ curl http://localhost:8080/metrics
...
# HELP http_requests_total Total http requests.
# TYPE http_requests_total counter
http_requests_total{code="200"} 40
http_requests_total{code="400"} 27
# HELP program_info Info about the program.
# TYPE program_info gauge
program_info{port="8080"} 1
...
```
