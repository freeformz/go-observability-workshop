# Prometheus

## Data Model

ALl data is stored as time series.

Metrics are identified by name AND any optional labels (key/value pairs)

Label names may contain ASCII letters, numbers, as well as underscores.

Label values may contain any Unicode characters.

```
<metric name>{<label name>=<label value>, ...}
http_requests_total{method="POST", handler="/messages", service="api"}
api_http_requests_total{method="POST", handler="/messages"}
```

Every unique label set is a new time series.

```
http_requests_total{method="POST", handler="/messages", service="api"}
http_requests_total{method="GET", handler="/messages", service="api"}
http_requests_total{method="GET", handler="/messages", service="foo"}
http_requests_total{method="POST", handler="/alerts", service="api"}
```

^4 different time series

High cardinality label values are bad. Don't use things like: ip addresses, user/account/record ids, etc.

## Naming

Metric names specify the general aspect of a system that is being measured.

Example: `http_requests_total` == the total number of http request received

*Should* have a single word prefix (aka namespace). Can be used to specify application or generic names.

Example: `api_frobs_total` == the total number of times the api frobs

*Must* have a single unit. Seconds or Bytes, Seconds not mixed with milliseconds.

*Should* use base units like: bytes (instead of milliseconds, nanoseconds, etc), bytes (instead of megabytes), etc

List of [base units](https://prometheus.io/docs/practices/naming/#base-units)

*Should* have a suffix describing the unit, in plural form. `total` is used for counts.

Examples:

* `http_request_duration_seconds`
* `node_memory_usage_bytes`
* `http_requests_total`
* `process_cpu_seconds_total`
* `foobar_build_info`

*Should* (really really really should) represent and measure the same thing across all label dimensions.

Examples:

* Request duration
* Bytes of data transfered
* Instantaneous resource usage as a percentage (CPU for example)

## General rule of thumb

Using the `sum()` and `avg()` functions over all dimensions of a given metric should be meaningful. If not, split the data up into multiple metrics. Example:

`queue_total{aspect="capacity", ...}` & `queue_total{aspect="depth", ...}` should be `queue_depth_total` and `queue_capacity_total` instead.

## Gotchas

* Avoid missing metrics: they break alerts and confuse rates. Initialize them at application startup.
* Have `total` and `failure` metrics instead of `success` and `failure` methods.

## Metric Types

https://prometheus.io/docs/concepts/metric_types/

### Counter

Cumulative metric representing a single, resetable, monotonically increasing value

Example usage: number of requests served, jobs run, errors, etc.

### Gauge

Metric representing a single numerical value. Usually a snapshot of some value at a specific time.

Example usage: speed of a car, temperature of a room, number of go routines, etc.

### Histogram

Composite metric used to sample observations, counting them into configurable buckets and sums.

Buckets are expressed as less than or equal an upper bounds target value. This makes histograms cumulative.

A histogram exposes multiple time series. Given a `<base name>`, the following time series are exposed:

* cumulative counters for the observation buckets, exposed as `<basename>_bucket{le="<uppser inclusive bound>"}`.
* the total sum of all observed values, exposed as `<basename>_sum`.
* the count of events that have been observed, exposed as `<basename>_count`

Example usage: request durations, response sizes, etc.

### Summary

Similar to histograms, but calculates quantiles over a sliding time window.

A summary exposes multiple time series. Given a `<base name>`, the following time series are exposed:

* streaming quantiles of observed evets, exposed as `<base name>{quantile="<q>"}`.
* the total sum of all observed values, exposed as `<base name>_sum`.
* the count of events that have ben observed as `<base name>_count`.

Example usage: request durations, response sizes, etc.

### Histogram vs Summary

Histogram:

* Configure buckets for expected ranges
* Observations are cheap (low impact on the observer)
* Server has to calculate quantiles (pXX)
* One time series per bucket + _sum + _count
* Quantile error is limited by the width of the relevant bucket
* Ad-hoc quantiles
* Ad-hoc aggregation

Summary:

* Pick quantiles and sliding window
* Observations are expensive (higher impact on the observer)
* Server doesn't have to calculate much for quantiles (pXX)
* One time series per quantile + _sum + _count
* Quantile error is limited by defined quantile objectives
* Preconfigured quantiles
* Generally not aggregatable

As a general rule, you should use histograms, but need to pay attention to you observation distribution to configure a suitable bucket layout.

^ Show example
