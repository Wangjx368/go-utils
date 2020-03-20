package trace

import (
	"github.com/astaxie/beego"
	"github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	enpoitUrl   = beego.AppConfig.String("zipkin_url")
	serviceName = beego.AppConfig.String("zipkin_service_name")
	ip          = beego.AppConfig.String("zipkin_service_ip")
)

type Trace struct {
	Tracer *zipkin.Tracer
	Span   zipkin.Span
}

func GetTracer() *zipkin.Tracer {
	// create a reporter to be used by the tracer
	reporter := httpreporter.NewReporter(enpoitUrl)

	// set-up the local endpoint for our service
	endpoint, _ := zipkin.NewEndpoint(serviceName, ip)

	// set-up our sampling strategy
	sampler := zipkin.NewModuloSampler(1)

	// initialize the tracer
	tracer, _ := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	return tracer
}
