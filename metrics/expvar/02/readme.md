# Expvar Metrics

`expvar` also includes the `Float` and `Int` types. Either one of these types can be used as either counters or gauges.

```go
c := expvar.NewInt("MyCounter") // 0
c.Add(1) // 1
c.Add(3) // 4

// OR

g := expvar.NewInt("MyGauge") // 0
go func() {
  for {
    g.Set(sampleSomethingFunction()) // set the value to whatever sampleSomethingFunc() returns
    time.Sleep(1 * time.Minute)
  }
}()
```

There are various bits of tooling in the Go ecosystem that can handle and work with `expvar` data.

One of those is `expvarmon`, which is useful for connecting to a server instrumented with `expvar` and observing it in real time.

## Exercise

Continuing from the last exercise, make a copy of the server and modify it to:

1. Expose a counter named `Requests` that counts the number of all requests processed by the server.
1. Expose a counter named `Errors` that counts the number of errors returned by the server.

Then use `hey` to generate traffic and `expvarmon` to monitor the effects of that traffic in real time.

Note: We'll assume that # of Successful Requests == `Requests` - `Errors`.

## Prerequisites

```console
$ gobin -u github.com/divan/expvarmon
...
```

## Give it a try

```console
$ go run server.go &
$ hey -c 1 -z 60m http://localhost:8080/ &
$ expvarmon -i 2s -ports="8080" -vars "mem:memstats.Alloc,mem:memstats.Sys,mem:memstats.HeapAlloc,
mem:memstats.HeapInuse,duration:memstats.PauseNs,duration:memstats.PauseTotalNs,
Requests,Errors,Port"
...
```
