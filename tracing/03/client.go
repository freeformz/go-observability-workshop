package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

const (
	apiURL = "http://localhost:8080/"
)

func apiGETRequest(ctx context.Context, url string) (*http.Request, error) {
	ctx, span := trace.StartSpan(ctx, "apiGETRequest")
	defer span.End()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return req.WithContext(ctx), nil
}

func processAPIResponse(ctx context.Context, r *http.Response) error {
	_, span := trace.StartSpan(ctx, "processAPIResponse")
	defer span.End()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "reading body")
	}

	log.Printf("%q\n", b)

	span.Annotate([]trace.Attribute{trace.Int64Attribute("response_bytes", int64(len(b)))}, string(b))
	return nil
}

func doAPIRequest(ctx context.Context, c *http.Client, url string) error {
	ctx, span := trace.StartSpan(ctx, "doAPIRequest")
	defer span.End()

	span.Annotate([]trace.Attribute{trace.StringAttribute("url", url)}, "api request to")

	req, err := apiGETRequest(ctx, url)
	if err != nil {
		return errors.Wrap(err, "creating api request")
	}

	res, err := c.Do(req)
	if err != nil {
		return errors.Wrap(err, "making api request")
	}
	defer res.Body.Close()

	return errors.Wrap(processAPIResponse(ctx, res), "processing api response")
}

func main() {
	je, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: "http://localhost:14268/api/traces",
		Process: jaeger.Process{
			ServiceName: "client",
			Tags: []jaeger.Tag{
				jaeger.StringTag("client_id", "workshop"),
				jaeger.BoolTag("client", true),
				jaeger.BoolTag("demo", true),
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create the Jaeger exporter: %v", err)
	}
	trace.RegisterExporter(je)                                            //register the exporter
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}) // demo, so always sample

	ctx := context.Background()
	var oct ochttp.Transport
	client := http.Client{Transport: &oct}

	for {
		ctx, span := trace.StartSpan(ctx, "top of loop")

		if err := doAPIRequest(ctx, &client, apiURL); err != nil {
			span.SetStatus(trace.Status{Code: trace.StatusCodeInternal, Message: err.Error()})
			log.Println("OOPS:", err)
		}

		span.End()

		time.Sleep(1 * time.Second)
	}
}
