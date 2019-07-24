# Micro Services Setup

Let's approximate a very simple service that is comprised of 2 services.
In a real world scenario there would be more than 2 services.

To create out two services, copy the prometheus/06 code to both `servicea/servicea.go` and `serviceb/serviceb.go`
Modify servicea to act like a proxy or API front-end to serviceb.
Make 25% of the requests that go to "/" on servica make a request to serviceb's "/slow".
Make the remaining requests to servicea's "/" make a request to serviceb's "/".
Proxy any bytes returned from serviceb through to the client making the original request using `io.Copy`.

Run hey against servicea "/".
Note: We'll be ignoring servicea's "/slow".

This is probably the most complicated step.
But if you are using microservices, this is is exemplary of what is happend, even if it's wrapped up in other packages.

When you are done, launch both services and use the `hey` command we've been using to send traffic to servicea's `/` route.

If everything is working you should see logs scroll by for both services.

Bonus activity: How do you correlate a request to servicea to the same request on serviceb?
Get/set/log an `X-Request-ID` header generated via the `github.com/google/uuid` package.
