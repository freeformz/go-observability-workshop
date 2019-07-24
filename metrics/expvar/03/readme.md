# Custom Expvar Metrics

expvar defines a `Var` interface. Anyting implementing that interface must return a valid JSON representation of the value to include in the `/debug/vars` output.
Because it's a JSON string, multiple values can be returned.

## Exercise

Bulding on Exercise #2, create a new type that returns the following values:

* The total number of requests (`Requests.Count`);
* The total time spent servicing requests (`Requests.Sum`);
* The average time spent per request (`Requests.Avg`).

Create a middleware that can be used to wrap the existing handler to do the actual instrumentation of the handler.

## Give it a try

```console
$ go run server.go &
$ hey -c 1 -z 60m http://localhost:8080/ &
$ expvarmon -i 2s -ports="8080" -vars "mem:memstats.Alloc,mem:memstats.Sys,mem:memstats.HeapAlloc,mem:memstats.HeapInuse,duration:memstats.PauseNs,duration:memstats.PauseTotalNs,Requests.Count,Errors,duration:Requests.Sum,duration:Requests.Avg"
...



hey -c 1 -z 60m http://localhost:8080/ &
$ expvarmon -i 2s -ports="8080" -vars "mem:memstats.Alloc,mem:memstats.Sys,mem:memstats.HeapAlloc,mem:memstats.HeapInuse,duration:memstats.PauseNs,duration:memstats.PauseTotalNs,Requests,Errors,Port"
...
```