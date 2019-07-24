# Intro

## What is Observability

Observability is an attribute of your system.

At its essence it is the ability to answer questions about your system.

Your system can be more, or less observable.

During this workshop we are going to focus on some of the different methods you can use to make your system(s) observable.

## Not Covered

We will not be covering the following topics:

1. Monitoring
1. Alerting
1. Setting up/maintaining prometheus
1. Setting up/maintaining jaeger

These are large topics that would basically require their own workshops.

## Before we begin

There are no wrong answers, just opportunities to learn something new.

This applies not only to you, the workshop attendees, but to me as well.

## 3 Pillars of Observability

1. Metrics
1. Logs
1. Traces

[![Observability Venn Diagram](../assets/obsrevability-ven.png)](https://peter.bourgon.org/blog/2017/02/21/metrics-tracing-and-logging.html)

## Metrics

Before we can talk about Metrics we need to talk about Measurements.

A measurement is a value pertaining to an aspect of your system at a given point in time.

A metric is one or more measurements aggregated using one or more statistical methods for a period of time.

Measurements are irregular (whey happen when they happen), while metrics are regular (they are emitted at regular intervals).

```text
┌───────────────────────────┐    ┌────────────────────────────┐
│ 1s worth of measurements  │ │  │     aggregated metric      │
└───────────────────────────┘    └────────────────────────────┘
  ┌──────────────┐            │
  │              │
  │   foo = 1    ├────────┐   │
  │              │        │         ┌────────────────────────┐
  └──────────────┘        │   │     │                        │
  ┌──────────────┐        │         │ foo = (                │
  │              │        │   │     │     avg: 2.75,         │
  │   foo = 5    │────────┤         │     max: 5,            │
  │              │        │   │     │     min: 0,            │
  └──────────────┘        ├────────▶│     sum: 9,            │
  ┌──────────────┐        │   │     │   count: 4,            │
  │              │        │         │    when: 1562367429,   │
  │   foo = 0    │────────┤   │     │      etc...            │
  │              │        │         │       )                │
  └──────────────┘        │   │     │                        │
  ┌──────────────┐        │         └────────────────────────┘
  │              │        │   │
  │   foo = 3    │────────┘
  │              │            │
  └──────────────┘
```

## Logs

Logs are some form of textual output about something happening inside of your system.

Everyone is probably familiar with logs.

Logs alert operators or other systems about what is happening in the system emitting the logs.

Opinions about the correct amount of verbosity, the audience for logs, and how they're formatted varies greatly.

Some believe logs should be reserved for messages destined only for an operator and should be human readable.

Others believe that logs should be structured and emitted in real time.

## Traces

A trace is a set of causally related events, triggered as a result of a logical operation, consolidated across the various components of an application.

A distributed trace contains events that cross process, network and security boundaries.

Traces allow engineers to understand the different services, systems, boundaries, and actions involved in the path of a request or other logical operation.

Forks or hops in execution flow are logical boundaries in a trace.

Most effective when every system or component involved in an operation is participates in tracing.

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

## Metrics - expvar

Go comes with it's own rudimentary metrics types in the `[expvar](https://golang.org/pkg/expvar/)` package.

`expvar` provides a "standardized interface to public variables" and "exposes these variables via HTTP".

[expvar exercises](../metrics/expvar)

### References

* [Monitoring Observability](https://medium.com/@copyconstruct/monitoring-and-observability-8417d1952e1c)
* [Monitoring and Observability — What’s the Difference and Why Does It Matter?](https://thenewstack.io/monitoring-and-observability-whats-the-difference-and-why-does-it-matter/)
* [Observability](https://en.wikipedia.org/wiki/Observability)
* [What is Observability](https://engineering.salesforce.com/what-is-observability-d175eb6cd2e4)
* [Monitoring And Observability](https://theagileadmin.com/2018/02/16/monitoring-and-observability/)
* [Introduction To Observability](https://docs.honeycomb.io/learning-about-observability/intro-to-observability/)
* [Distributed Systems Observability](https://www.oreilly.com/library/view/distributed-systems-observability/9781492033431/ch01.html)
