package tracer

import (
	"fmt"
	"net/http"
	"time"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func GetTraceableHTTPClient(timeout *time.Duration, resourceName string) *http.Client {
	var client *http.Client
	if timeout != nil {
		client = &http.Client{
			Timeout: *timeout,
		}
	} else {
		client = &http.Client{}
	}
	namer := func(r *http.Request) string {
		return fmt.Sprintf("ext_%s", resourceName)
	}

	opt := httptrace.RTWithResourceNamer(namer)
	client = httptrace.WrapClient(client, opt)
	return client
}
