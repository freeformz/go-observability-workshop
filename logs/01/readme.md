# Go stdlib Log Package

Go's [stdlib log package](https://golang.org/pkg/log/).

Provides:

* a default instance as a global.
* a `Logger` type.

`Logger` values provide:

* A means of configuring a log message prefix via flags and prefix functions;
* A means to just log a message via `Print{,f,ln}`;
* A means of logging and panicing via `Panic{,f,ln}`;
* A means of logging and exiting the program via `Fatal{,f,ln}`;
* A means of replacing the log destination with any io.Writer;

The last is useful when testing.

[![tweet about testing logs](../../assets/twitter3.png)](https://twitter.com/peterbourgon/status/1151886045861928961)

## Exercise

Let's create a simple app that has the following characteristics:

* Listens locally on a port for HTTP requests;
* HTTP requests to `/` do some "work" that take between 1 and 100 ms to complete;
* If the work is successful, the app responds with `:-)`;
* The type of "work" the app is doing errors 25% of the time;
* When the work errors the app responds with `:-(`;
* When the work errors an error that says `OMG Error!` is logged;
* Logs the complete url to connect to on startup;
* Logs any errors returned by `http.ListenAndServe`;
* Each request is logged, allong with the HTTP method, the url requests, the return code and the duration of the request in seconds;
* Any logged output contains the filename and location of that emitted the log lines;
* Utilizes Go's stdlib log package for all logging.

Bonus ativity if you are done quickly: Do it without using the global log instance.

## Prerequisites

```console
$ GO111MODULE=off go get -u github.com/myitcv/gobin    # Used to install useful binaries
$ gobin -u github.com/rakyll/hey                       # Used to send test traffic to our server(s)
$ go get -u github.com/pkg/errors                      # Use this at least until go1.13's new errors package
...
```

## Give it a try

```console
$ go run logs/01/server.go &
$ hey -c 1 -z 60m http://localhost:8080/
2019/07/22 13:30:35 01.go:51: Listening at: http://localhost:8080
...
2019/07/22 13:30:44 01.go:23: Work took 0.024s
2019/07/22 13:30:44 01.go:35: Error: OMG Error!
2019/07/22 13:30:44 01.go:29: GET "/" => 400 (0.028s)
2019/07/22 13:30:44 01.go:23: Work took 0.079s
2019/07/22 13:30:44 01.go:29: GET "/" => 200 (0.081s)
...
```
