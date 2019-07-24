# Vectored Histograms

## Breakdown

You can think of histograms as a bar chart, where each bar in the chart accumulates the values that fall within a range.

```text
        ▲
   90  ─┤
        │
   80  ─┤
R       │
e  70  ─┤                     ┌──────┐
q       │                     │      │
u  60  ─┤              ┌──────┤      ├──────┐
e       │              │      │      │      │
s  50  ─┤       ┌──────┤      │      │      ├──────┐
t       │       │      │      │      │      │      │
   40  ─┤       │      │      │      │      │      │
C       │┌──────┤      │      │      │      │      ├──────┐
o  30  ─┤│      │      │      │      │      │      │      │
u       ││      │      │      │      │      │      │      │
n  20  ─┤│      │      │      │      │      │      │      │
t       ││      │      │      │      │      │      │      │
   10  ─┤│ 0-9  │10-19 │20-29 │30-39 │40-49 │50-59 │60-69 │
        ││  ms  │  ms  │  ms  │  ms  │  ms  │  ms  │  ms  │
        │└──────┴──────┴──────┴──────┴──────┴──────┴──────┘
         ──────────────────────────────────────────────────▶
                          Request Duration
```

This is also an example of a "normal" distribution.

* 35 requests in the 0-9ms bucket
* 50 requests in the 10-19ms bucket
* 60 requests in the 20-29ms bucket
* 70 requests in the 30-39ms bucket
* 60 requests in the 40-49ms bucket
* 50 requests in the 50-59ms bucket
* 35 requests in the 60-69ms bucket

Or about 360 requests total.

What is the 90p quantile of this?

360 * .9 = 324

35 (0-9) + 50 (10-19) + 60 (20-29) + 70 (30-39) + 60 (40-49) + 50 (50-59) = 325

So, the p90 of this histogram is somewhere in the 50-59ms range.
Prometheus uses linear interpolation and would likely calculate ~59ms or so.

Note: Linear interpolation works best for larger (i.e. high traffic sites) sample sizes, because "math".

Let's pretend that the graph above represents the distribution of 1 minute worth of request traffic.

## Exercise

Convert the `CounterVec` to a `HistogramVec` so we can also sample the request durations.

Bonus activity: Explain why we no longer need a counter.

## Give it a try

```console
$ go run 04.go &
$ hey -c 1 -z 60m http://localhost:8080/ &
$ curl http://localhost:8080/metrics
...
# HELP http_request_duration_seconds HTTP request duration.
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{code="200",le="0.01"} 0
http_request_duration_seconds_bucket{code="200",le="0.02"} 0
http_request_duration_seconds_bucket{code="200",le="0.03"} 0
http_request_duration_seconds_bucket{code="200",le="0.04"} 6
http_request_duration_seconds_bucket{code="200",le="0.05"} 9
http_request_duration_seconds_bucket{code="200",le="0.06"} 13
http_request_duration_seconds_bucket{code="200",le="0.07"} 21
http_request_duration_seconds_bucket{code="200",le="0.08"} 22
http_request_duration_seconds_bucket{code="200",le="0.09"} 32
http_request_duration_seconds_bucket{code="200",le="+Inf"} 38
http_request_duration_seconds_sum{code="200"} 2.5563333900000003
http_request_duration_seconds_count{code="200"} 38
http_request_duration_seconds_bucket{code="400",le="0.01"} 6
http_request_duration_seconds_bucket{code="400",le="0.02"} 15
http_request_duration_seconds_bucket{code="400",le="0.03"} 27
http_request_duration_seconds_bucket{code="400",le="0.04"} 27
http_request_duration_seconds_bucket{code="400",le="0.05"} 27
http_request_duration_seconds_bucket{code="400",le="0.06"} 27
http_request_duration_seconds_bucket{code="400",le="0.07"} 27
http_request_duration_seconds_bucket{code="400",le="0.08"} 27
http_request_duration_seconds_bucket{code="400",le="0.09"} 27
http_request_duration_seconds_bucket{code="400",le="+Inf"} 27
http_request_duration_seconds_sum{code="400"} 0.46929777899999997
http_request_duration_seconds_count{code="400"} 27
# HELP program_info Info about the program.
# TYPE program_info gauge
program_info{port="8080"} 1
...
```
