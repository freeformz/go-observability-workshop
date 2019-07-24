# Logs & Logging

We won't be talking about logs in the datbase sense (Ex: Write Ahead Logs).

We're talking about logs that your program outputs either to stdout, to a file, and/or to an external system.

Logs provide insight into what is happening in your system in near real time.

Generally logs are timestamped, either by the system that is creating them, or by the systems that collect and/or route them.

## Structured Logs

There are generally 2 views about logging and logs.

One view argues that logs should be *only* for humans.

The other argues that logs should be for machines.

I've personally believed both views at different times in my career and currently emphatically believe that: It depends.

What does it depend on?

If the *only* observability tool you have is logs, then yes, log *everything*, *all the time*, and make the logs *structured* so they are easy to parse by other systems (like Splunk, elastic search, Papertrail, etc).
You'll end up using and/or writing tools to do so.
This is historically how we've done things @ Heroku.
We even have systems and software that turn logs into metrics.

But if you have more robust observability tools, reduce your logging to important events meant largely for human consumption. This is especially useful when combined with metrics and tracing.

See also [this twitter thread](https://twitter.com/VladimirVivien/status/1151899814076043264?s=20):

![twitter convo](../../assets/twitter1.png)

[Dave Cheney](https://dave.cheney.net/2015/11/05/lets-talk-about-logging):

    I believe that there are only two things you should log:

        Things that developers care about when they are developing or debugging software.
        Things that users care about when using your software.

