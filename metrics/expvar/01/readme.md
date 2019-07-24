# Expvar

Just like with logging, Go's stdlib includes a basic metrics package named `expvar`.

Just importing `expvar` exposes anything registered with it via the `http.DefaultServeMux` at `/debug/vars`.

`expvar` includes types that can act as counters (`*expvar.Int`), as well as expose program information (`*expvar.String`).

Note: The expvar pacakge works primarily with package level registration.

```go
v := expvar.NewString("Value")
v.Set("my value")
```

Before we dive deeper into metrics...

## Exercise

Starting with a copy of the logs/02 server:

1. Modify the server to expose the listening port via an `*expvar.String`;
1. Use curl & jq to look at this value and what else expvar exports.

## Prerequisites

[`jq`](https://stedolan.github.io/jq/) if you want to pretty print or query any json.
[`curl`](https://curl.haxx.se/) to make http calls against the server.

## Give it a try

```console
$ go run server.go &
$ curl -s http://localhost:8080/debug/vars | jq .Port
"8080"
$ curl -s http://localhost:8080/debug/vars | jq -c
{"Port":"8080","cmdline":["/var/folders/f7/r5gtrkh53nl49cntlpmhl_s5rz0nzb/T/go-build319404131/b001/exe/server"],"memstats":{"Alloc":449016,"TotalAlloc":449016,"Sys":70453248,"Lookups":0,"Mallocs":1579,"Frees":119,"HeapAlloc":449016,"HeapSys":66715648,"HeapIdle":65265664,"HeapInuse":1449984,"HeapReleased":0,"HeapObjects":1460,"StackInuse":393216,"StackSys":393216,"MSpanInuse":22176,"MSpanSys":32768,"MCacheInuse":13888,"MCacheSys":16384,"BuckHashSys":2607,"GCSys":2240512,"OtherSys":1052113,"NextGC":4473924,"LastGC":0,"PauseTotalNs":0,"PauseNs":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"PauseEnd":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"NumGC":0,"NumForcedGC":0,"GCCPUFraction":0,"EnableGC":true,"DebugGC":false,"BySize":[{"Size":0,"Mallocs":0,"Frees":0},{"Size":8,"Mallocs":45,"Frees":0},{"Size":16,"Mallocs":609,"Frees":0},{"Size":32,"Mallocs":104,"Frees":0},{"Size":48,"Mallocs":193,"Frees":0},{"Size":64,"Mallocs":94,"Frees":0},{"Size":80,"Mallocs":26,"Frees":0},{"Size":96,"Mallocs":50,"Frees":0},{"Size":112,"Mallocs":19,"Frees":0},{"Size":128,"Mallocs":24,"Frees":0},{"Size":144,"Mallocs":14,"Frees":0},{"Size":160,"Mallocs":23,"Frees":0},{"Size":176,"Mallocs":5,"Frees":0},{"Size":192,"Mallocs":5,"Frees":0},{"Size":208,"Mallocs":23,"Frees":0},{"Size":224,"Mallocs":12,"Frees":0},{"Size":240,"Mallocs":0,"Frees":0},{"Size":256,"Mallocs":25,"Frees":0},{"Size":288,"Mallocs":13,"Frees":0},{"Size":320,"Mallocs":2,"Frees":0},{"Size":352,"Mallocs":33,"Frees":0},{"Size":384,"Mallocs":27,"Frees":0},{"Size":416,"Mallocs":5,"Frees":0},{"Size":448,"Mallocs":3,"Frees":0},{"Size":480,"Mallocs":3,"Frees":0},{"Size":512,"Mallocs":0,"Frees":0},{"Size":576,"Mallocs":6,"Frees":0},{"Size":640,"Mallocs":4,"Frees":0},{"Size":704,"Mallocs":2,"Frees":0},{"Size":768,"Mallocs":0,"Frees":0},{"Size":896,"Mallocs":5,"Frees":0},{"Size":1024,"Mallocs":17,"Frees":0},{"Size":1152,"Mallocs":4,"Frees":0},{"Size":1280,"Mallocs":1,"Frees":0},{"Size":1408,"Mallocs":1,"Frees":0},{"Size":1536,"Mallocs":0,"Frees":0},{"Size":1792,"Mallocs":7,"Frees":0},{"Size":2048,"Mallocs":3,"Frees":0},{"Size":2304,"Mallocs":3,"Frees":0},{"Size":2688,"Mallocs":2,"Frees":0},{"Size":3072,"Mallocs":2,"Frees":0},{"Size":3200,"Mallocs":0,"Frees":0},{"Size":3456,"Mallocs":0,"Frees":0},{"Size":4096,"Mallocs":24,"Frees":0},{"Size":4864,"Mallocs":1,"Frees":0},{"Size":5376,"Mallocs":1,"Frees":0},{"Size":6144,"Mallocs":7,"Frees":0},{"Size":6528,"Mallocs":0,"Frees":0},{"Size":6784,"Mallocs":0,"Frees":0},{"Size":6912,"Mallocs":0,"Frees":0},{"Size":8192,"Mallocs":4,"Frees":0},{"Size":9472,"Mallocs":8,"Frees":0},{"Size":9728,"Mallocs":0,"Frees":0},{"Size":10240,"Mallocs":0,"Frees":0},{"Size":10880,"Mallocs":0,"Frees":0},{"Size":12288,"Mallocs":0,"Frees":0},{"Size":13568,"Mallocs":0,"Frees":0},{"Size":14336,"Mallocs":0,"Frees":0},{"Size":16384,"Mallocs":0,"Frees":0},{"Size":18432,"Mallocs":0,"Frees":0},{"Size":19072,"Mallocs":0,"Frees":0}]}}
```
