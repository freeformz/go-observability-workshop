# Enter Distributed Tracing

Distributed tracing takes a request centric view and captures detailed information about the request as it flows through a service and out to others.
Instrumentors (devs w/libraries and/or infrastructure) attach contextual metadata to each request and ensure this metadata is passed around during execution between services.
Instrumentors annotate the traces with relevant information (events, metrics, metadata, etc) at various points.
Explicit metadata is added to enforce causal references between systems and prior events (think a http request that leads to a background job in a queue).

Let's refer back to the diagram in the intro....

```text
          Time
    ─────────────────────────────────────────────────────────────────────────────────────────────▶
  │ ┌────────────────────────────────────────────────────────────────────────────────────────────┐
  │ │<-------------------------------- Start to finish 320ms ----------------------------------->│
  │ │                                                                                            │
  │ │                                                                                            │
C │ └────────────────────────────────────────────────────────────────────────────────────────────┘
o │ ┌──────────────┐                                                              ┌──────────────┐
n │ │A <--50ms--> B│                                                              │B' <-50ms-> A'│
c │ │              │                                                              │              │
u │ │              │                                                              │              │
r │ └──────────────┘                                                              └──────────────┘
r │                 ┌─────────────────────────────────────────────────────────────┐
e │                 │B <------------------------ 250 ms ----------------------> B'│
n │                 │                                                             │
c │                 │                                                             │
y │                 └─────────────────────────────────────────────────────────────┘
  │                 ┌─────────────┐
  │                 │<---20ms---> │
  │                 │             │
  │                 │    Redis    │
  │                 │ Transaction │
  │                 └─────────────┘
  │                 ┌─────────────────────────────────────────────────────┐
  │                 │B <------------------- 230 ms --------------------> C│
  │                 │                                                     │
  │                 │                                                     │
  │                 └─────────────────────────────────────────────────────┘
  │                       ┌───────────────────────────────────────────────┐
  │                       │C <----------------- 210 ms ----------------> D│
  │                       │                                               │
  │                       │                                               │
  ▼                       └───────────────────────────────────────────────┘
```

The core building block of distributed tracing is the `Span`.
A `Span` has a`TraceID`. The `TraceID` is globally unique and every span with the `TraceID` belongs to the same trace.

Each `Span` has it's own id and may have a `ParentID`. Spans without `ParentID`s are called root spans, while spans with a `ParentID` are called child spans. Multiple spans can have the same `ParentID`, making them all children of the span with that id. Child spans can themselves have child spans.

Additionally, spans have:

* a `Status` field, think of them as similar to http response codes;
* tags that provide contextual key-value information about the span. In that way are similar to tags in prometheus; And some traceing implmentations (like opencensus) will use them when generating metrics;
* annotations that work like structured logs that describing something that happened at a given point in the lifecycle of the span.
* a start time
* and end time

Spans may also be `linked` to other spans to form causal relationships. This is often used to link the spans of a batch job to the spans of a http request that caused the job to be queued.

Because spans can contain a lot of data, they are often sampled using different methods. Today, we'll be sampling *all* spans, but in production spans are generally sampled at less than full fidelity.

Span information can also be "propagated" from one system to another in a few different formats. Open census, Zipkin, and other systems use the `b3` propagation format: https://github.com/openzipkin/b3-propagation.

Review:  https://godoc.org/go.opencensus.io/trace

Key apis:

```go
// Exporting
trace.RegisterExporter(exporter)

// Sampling
trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

// Creating a new span
ctx, span := trace.StartSpan(ctx, "example.com/Run")
defer span.End()

// Annotate the span with structured log like data.
span.Annotate(
  []trace.Attribute(
    trace.StringAttribute("thing","description"),
    trace.BoolAttribute("yes",true),
  ),
  "message",
)

// The jaeger config we'll be using
je, err := jaeger.NewExporter(jaeger.Options{
  CollectorEndpoint: "http://localhost:14268/api/traces",
  Process: jaeger.Process{
    ServiceName: "serviceName",
    Tags: []jaeger.Tag{
      jaeger.StringTag("a_tag", "a_value"),
      //...
    },
  },
})
if err != nil {
  log.Fatalf("Failed to create the Jaeger exporter: %v", err)
}

// Use the opencensus transport to propigate trace/span ids in http requests
var oct ochttp.Transport
c := http.Client{Transport: &oct, Timeout: 2 * time.Second} 

// net/http Middleware
// https://godoc.org/go.opencensus.io/plugin/ochttp#WithRouteTag
ochttp.WithRouteTag(myHandler,"/")

// ochttp.Handler wraps a http.Handler ensuring the server can understand the b3 propigation format.
// https://godoc.org/go.opencensus.io/plugin/ochttp/propagation/b3
// https://godoc.org/go.opencensus.io/plugin/ochttp#Handler
var pf b3.HTTPFormat
&ochttp.Handler{
  Handler:     mux,
  Propagation: &pf,
}
```

Future Reading:

* https://www.w3.org/TR/trace-context/ - This specification defines standard headers and value format to propagate context information that enables distributed tracing scenarios. The specification standardizes how context information is sent and modified between services. Context information uniquely identifies individual requests in a distributed system and also defines a means to add and propagate provider-specific context information.

## Exercise

Move servicea and serviceb from logging based tracing to opencensus tracing using spans.
The key APIs you need to use are listed above.

Let's go through my version of the exercise first.

When you are done untar jaeger's `all-in-one` process and run it. It will use an in-memory data store for the traces:

```console
$ ./jaeger-all-in-one &
$ go run servicea/servicea.go &
$ go run serviceb/serviceb.go &
$ hey -c 2 -z 60m http://localhost:8080 &
$ open http://localhost:16686   # To view the jaeger ui
...
```

Things to try in the UI:

* Search for tag `http.status_code=400`
* explore the "Operation" drop down
* click into spans with and without errors.
* expand all span details.

