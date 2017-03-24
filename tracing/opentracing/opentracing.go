package opentracing

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"google.golang.org/grpc"
)

type Collector struct {
	URL     string
	Options []zipkin.HTTPOption
	Client  zipkin.Collector
}

type Recorder struct {
	Debug       bool
	Host        string
	ServiceName string
	Options     []zipkin.RecorderOption
	Client      zipkin.SpanRecorder
}

type Tracer struct {
	Collector Collector
	Recorder  Recorder
	Client    opentracing.Tracer
}

func NewTracer(t Tracer) (*Tracer, error) {
	collector, err := zipkin.NewHTTPCollector(t.Collector.URL, t.Collector.Options...)
	if err != nil {
		return nil, err
	}
	recorder := zipkin.NewRecorder(collector, t.Recorder.Debug, t.Recorder.Host, t.Recorder.ServiceName)
	zktracer, err := zipkin.NewTracer(recorder)
	if err != nil {
		return nil, err
	}

	t.Collector.Client = collector
	t.Recorder.Client = recorder
	t.Client = zktracer

	return &t, nil
}

func (t *Tracer) ClientInterceptor() grpc.UnaryClientInterceptor {
	return otgrpc.OpenTracingClientInterceptor(t.Client)
}

func (t *Tracer) ServerInterceptor() grpc.UnaryServerInterceptor {
	return otgrpc.OpenTracingServerInterceptor(t.Client)
}

func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.Client.StartSpan(operationName, opts...)
}

func (t *Tracer) StartChildSpan(operationName string, parent opentracing.Span) opentracing.Span {
	return t.Client.StartSpan(operationName, opentracing.ChildOf(parent.Context()))
}

func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return t.Client.Inject(sm, format, carrier)
}

func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return t.Client.Extract(format, carrier)
}
