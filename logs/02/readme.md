# Structured Logs / Log Levels

There are generally 2 views about logging and logs.

One view argues that logs should be *only* for humans.

The other argues that logs should be for machines.

I've personally believed both views at different times in my career and currently emphatically believe that: It depends.

What does it depend on?

If the only observability tool you have is logs, then yes, log *everything*, *all the time*, and make the logs *structured* so they are easy to parse by other systems (like Splunk, elastic search, Papertrail, etc). You'll end up using and/or writing tools to do so. This is historically how we do things @ Heroku.

But if you have more robust observability tools, reduce your logging to important events meant largely for human consumption. This is especially useful when combined with metrics and tracing.

See also [this twitter thread](https://twitter.com/VladimirVivien/status/1151899814076043264?s=20):

![twitter convo](../../assets/twitter1.png)

[Dave Cheney](https://dave.cheney.net/2015/11/05/lets-talk-about-logging):

    I believe that there are only two things you should log:

        Things that developers care about when they are developing or debugging software.
        Things that users care about when using your software.

## Log Levels

For a *long* time the only structured bit of information that you could add to a log message was it's Severity Level.
The levels, encoded in various syslog RFCS (5424, 3164) are part of syslog's [`PRIVAL`](https://tools.ietf.org/html/rfc5424#section-6.2.1)) field and are as follows:

| Value | Severity  | Keyword | Description                       | Condition                                                                       |
|-------|-----------|---------|-----------------------------------|---------------------------------------------------------------------------------|
| 0     | Emergency | emerg   | System is unusable                | A panic condition                                                               |
| 1     | Alert     | alert   | Action must be taken immediately  | A condition that should be corrected immediately                                |
| 2     | Critical  | crit    | Critical condition                | Hard device errors                                                              |
| 3     | Error     | err     | Error conditions                  |                                                                                 |
| 4     | Warning   | warning | Warning conditions                |                                                                                 |
| 5     | Notice    | notice  | Normal but significant conditions | Non error condition that may require special handling                           |
| 6     | Info      | info    | Informational messages            |                                                                                 |
| 7     | Debug     | debug   | Debug-level messages              | Messages that contain information normally of se only when debugging a program. |

## Structured Logs

Previously our program was outputting logs that look something like this:

```text
2019/07/17 18:23:54 GET "/foo" => 200 (0.028s)
2019/07/17 18:23:54 OMG Error!
```

That's not too hard to parse for either a human or machine, but formats like JSON or logfmt are easier for machines to parse and extract useful information out of.
A format like this is way more machine parsable, and honestly becomes pretty easy to read once you are used to it:

```text
time="2019-07-17T21:21:16-07:00" level=info duration=0.06609306 method=GET path=/foo status=200
time="2019-07-17T21:21:16-07:00" level=error msg="OMG Error!" method=GET path=/foo
```

The format of ^ is called [logfmt](https://brandur.org/logfmt).

Another common option is JSON formatted logs.

```json
{"level":"info","msg":"Listening at: http://localhost:8080","time":"2019-07-17T21:26:08-07:00"}
{"duration":0.032404677,"level":"info","method":"GET","msg":"","path":"/foo","status":200,"time":"2019-07-17T21:26:21-07:00"}
{"level":"error","method":"GET","msg":"OMG Error!","path":"/foo","time":"2019-07-17T21:26:21-07:00"}
```

JSON is used in a lot of places and enables stuff like this:

```console
$ cat logs.json | jq  'select( .level == "error")
{
  "level": "error",
  "method": "GET",
  "msg": "OMG Error!",
  "path": "/foo",
  "time": "2019-07-17T21:30:31-07:00"
}
{
  "level": "error",
  "method": "GET",
  "msg": "OMG Error!",
  "path": "/foo",
  "time": "2019-07-17T21:30:31-07:00"
}
```

## logrus

There are several logging packages out there that provide more features like structured logs, filters, hooks, etc.

These features are important as your program's operational needs grow.

The one I'm most familiar with is [`github.com/sirupsen/logrus`](https://github.com/sirupsen/logrus), so we'll use that package to
convert the unstructured logging in our server to structured logs.

`github.com/sirupsen/logrus` provides:

* Structured logging in multiple formats
* Colored logging
* Logging Hooks
* Logging Levels
* Log rotation support
* Fatal Handlers (triggered when any fatal level messge is logged)

One example of using the "Hooks" feature is the package [`rollrus`](https://github.com/heroku/rollrus), which most Go programs @ Heroku use to report `Panic` and `Fatal` log calls to [Rollbar](https://rollbar.com/).
Rollbar collects, depuplicates and otherwise helps make sense of the errors coming from your system(s).
There are other services out there that provide similar services, but Rollbar is the one we currently use @ Heroku.

Logrus also allows for fields to be "curried" into a logger that will apply those fields to each log emitted.
This is my generally recommended approach.

```go
  log = log.WithField("app","logs-02-server")
  log.WithField("err",err).Error("OMG Error!")
```

## Exercise

Continuing from logs/01...

Use `github.com/sirupsen/logrus` to accomplish 2 things:

1. Add structure to the logs being emitted;
1. Add log levels to the logs;
1. Create a curried logger containing a field named `app` with a value of `logs-02-server`;
1. Use the curried logger for all logging;

Bonus ativity if you are done quickly: Do it without using the global logrus instance.

## Prerequisites

```console
$ go get -u github.com/sirupsen/logrus
...
```

## Give it a try

```console
$ go run server.go &
$ hey -c 1 -z 60m http://localhost:8080/
time="2019-07-22T14:12:44-07:00" level=info msg="Listening at: http://localhost:8080" app=logs-02-server
...
time="2019-07-22T14:12:47-07:00" level=info msg="Work complete" app=logs-02-server method=GET path=/ work_seconds=0.024
time="2019-07-22T14:12:47-07:00" level=error msg="OMG Error!" app=logs-02-server method=GET path=/
time="2019-07-22T14:12:47-07:00" level=info app=logs-02-server duration=0.026662857 method=GET path=/ status=400
time="2019-07-22T14:12:47-07:00" level=info msg="Work complete" app=logs-02-server method=GET path=/ work_seconds=0.079
time="2019-07-22T14:12:47-07:00" level=info app=logs-02-server duration=0.079135422 method=GET path=/ status=200
time="2019-07-22T14:12:48-07:00" level=info msg="Work complete" app=logs-02-server method=GET path=/ work_seconds=0.021
...

OR

{"app":"logs-02-server","level":"info","msg":"Listening at: http://localhost:8080","time":"2019-07-22T14:13:33-07:00"}
...
{"app":"logs-02-server","level":"info","method":"GET","msg":"Work complete","path":"/","time":"2019-07-22T14:13:34-07:00","work_seconds":0.024}
{"app":"logs-02-server","level":"error","method":"GET","msg":"OMG Error!","path":"/","time":"2019-07-22T14:13:34-07:00"}
{"app":"logs-02-server","duration":0.026852646,"level":"info","method":"GET","msg":"","path":"/","status":400,"time":"2019-07-22T14:13:34-07:00"}
{"app":"logs-02-server","level":"info","method":"GET","msg":"Work complete","path":"/","time":"2019-07-22T14:13:34-07:00","work_seconds":0.079}
{"app":"logs-02-server","duration":0.082836129,"level":"info","method":"GET","msg":"","path":"/","status":200,"time":"2019-07-22T14:13:34-07:00"}
{"app":"logs-02-server","level":"info","method":"GET","msg":"Work complete","path":"/","time":"2019-07-22T14:13:34-07:00","work_seconds":0.021}
...
```
