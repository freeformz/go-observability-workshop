# Exposing Application Prometheus Metrics

## Breakdown

In the Prometheus [intro](../intro.md) we learned about the different types of metrics that Prometheus supports:

* Counters
* Gauges
* Histograms
* Summaries

The `Counter` and `Gauge` types are interface that work very much like the expvar types.

A `Counter` can be `Add()` or `Inc()` (incremented; short for `Add(1)`).

A `Guage` can be `Set()`, `Inc()`, `Dec()`, `Add()`, or `Sub()`.

## Exercise

Convert the expvar `Int`s to prometheus `Counter`s.
Leave the expvar `Port` for now, we'll cover that in the next section.

When done, use curl to view '/metrics' and `hey` to induce load.

## Give it a try

```console
$ go run server.go &
$ hey -c 1 -z 60m http://localhost:8080/ &
$ curl http://localhost:8080/metrics
...
# HELP http_errors_total Total http errors.
# TYPE http_errors_total counter
http_errors_total 54
# HELP http_requests_total Total http requests.
# TYPE http_requests_total counter
http_requests_total 200
...
```
