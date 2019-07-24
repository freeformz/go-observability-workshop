# Scraping Metrics

Prometheus has a few different modes of operation:

* Pulling (Scraping)
* Pushing

Pushing metrics in Prometheus is generally done via the Push Gateway, which is a seperate service.
It's generally recommended that pushing metrics not be used. So we won't be covering it here.

The common and recommended way to use prometheus it to have the prometheus scrape metrics from the various services that need to be monitored.
A Prometheus instance can also pull from another prometheus instance, allowing a federated and/or tiered deployment model.
This is a topic for another workshop though.

But it is important to get a feeling for how how exported metrics are represented in prometheus.

## Exercise

Download (if you don't already have it) and untar/gzip the [prometheus binaries](https://prometheus.io/download/) for your operating system to the 05 directory.

In the end you should have a directory named something like: `prometheus-<version>.<OS>-amd64`.

The provided prometheus.yaml file scrapes the local service every 10s. The generally recommended scrape interval is between 10 and 60s.

Start up Prometheus using the provided prometheus.yml file so that it scrapes samples from our server.

Starting with the code from the last exercise, modify it so that requests to `/slow` respond between 100ms and 300ms, randomly determined, without erroring. Add a corresponding label for the different handler.

Run 2 instaces of `hey -c 1 -z 60m http://localhost:8080/` & `hey -c 1 -z 60m http://localhost:8080/slow` against your app to generate request load.

Use the prometheus dashboard (`http://localhost:9090`) to experiment with the following queries:

Note: Use the "Console" to start, then take a look at the "Graph" tab

```text
http_request_duration_seconds_count
^ count of all HTTP requests

http_request_duration_seconds_count{code="200"}
^ count of "HTTP OK" requests

rate(http_request_duration_seconds_count{code="200"}[1m])
^ avg per/second "HTTP OK" requests over the last minute

rate(http_request_duration_seconds_count{code="200"}[5m])
^ avg per/second "HTTP OK" requests over the last 5 minutes.
```

^ Compare the last two

```text
rate(http_request_duration_seconds_count[5m])[30m:1m]
^ 5-minute avg per/second rate of http requests for the past 30 minutes, with a resolution of 1 minute.

sum by(code) (rate(http_request_duration_seconds_count[1m]))
^ avg per/second HTTP requests by code over the last minute

sum by(code) (rate(http_request_duration_seconds_count[5m]))
^ avg per/second HTTP requests by code over the last 5 minutes
```

```text
http_request_duration_seconds_bucket
^ count of HTTP requests per histogram bucket

http_request_duration_seconds_bucket{code="200"}
^ count of "HTTP OK" requests per histogram bucket per handler

http_request_duration_seconds_bucket{code="200"}[1m]
^ (range vector) 1 value per `evaluation_interval` (see config)
count of "HTTP OK" requests per histogram bucket per handler for the last minute
```

```text
histogram_quantile(0.90, http_request_duration_seconds_bucket{code="200"})
^ p90 of "HTTP OK" requests per handler cumulatively

histogram_quantile(0.90, rate(http_request_duration_seconds_bucket{code="200"}[1m]))
^ p90 of "HTTP OK" requests per handler for the last minute

histogram_quantile(0.90, rate(http_request_duration_seconds_bucket{code="200"}[5m]))
^ p90 of "HTTP OK" requests per handler for the last 5 minutes

histogram_quantile(0.90, sum(rate(http_request_duration_seconds_bucket[1m])) by (code, le))
^ p90 of HTTP requests by code for the last minute

histogram_quantile(0.90, sum(rate(http_request_duration_seconds_bucket[1m])) by (handler, le))
^ p90 of HTTP requests by handler for the last minute

histogram_quantile(0.90, sum(rate(http_request_duration_seconds_bucket[1m])) by (code, handler, le))
^ p90 of HTTP requests by (code,handler) for the last minute
```

## Prerequisites

```console
$ curl -O -L https://github.com/prometheus/prometheus/releases/download/v2.11.1/prometheus-2.11.1.darwin-amd64.tar.gz
OR
$ curl -O -L https://github.com/prometheus/prometheus/releases/download/v2.11.1/prometheus-2.11.1.linux-amd64.tar.gz
OR
$ curl -O -L https://github.com/prometheus/prometheus/releases/download/v2.11.1/prometheus-2.11.1.windows-amd64.tar.gz
$ tar zxf prometheus-2.11.1.darwin-amd64.tar.gz # or linx || windows
$ cd prometheus-2.11.1*
$ ./prometheus --config.file="../prometheus.yaml" &
$ go run server.go &
$ open http://localhost:9090
```
