package opentracing

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"google.golang.org/grpc"
)

// Collector 呈現了一個追蹤器內的收集器。
type Collector struct {
	URL     string
	Options []zipkin.HTTPOption
	Client  zipkin.Collector
}

// Recorder 呈現了一個追蹤器內的記錄器。
type Recorder struct {
	Debug       bool
	Host        string
	ServiceName string
	Options     []zipkin.RecorderOption
	Client      zipkin.SpanRecorder
}

// Tracer 呈現了一個追蹤器的實體。
type Tracer struct {
	Collector Collector
	Recorder  Recorder
	Client    opentracing.Tracer
}

// NewTracer 會接收一個只有設定的追蹤器建構體，並以此來建立一個新的、可用實體追蹤器。
func NewTracer(t Tracer) (*Tracer, error) {
	// 建立收集器。
	collector, err := zipkin.NewHTTPCollector(t.Collector.URL, t.Collector.Options...)
	if err != nil {
		return nil, err
	}
	// 建立記錄器。
	recorder := zipkin.NewRecorder(collector, t.Recorder.Debug, t.Recorder.Host, t.Recorder.ServiceName)
	// 建立追蹤器。
	zktracer, err := zipkin.NewTracer(recorder)
	if err != nil {
		return nil, err
	}
	// 將這些客戶端存放到追蹤器的建構體內。
	t.Collector.Client = collector
	t.Recorder.Client = recorder
	t.Client = zktracer

	return &t, nil
}

// ClientInterceptor 會回傳一個可供 gRPC 中 Dial 函式使用的 Interceptor。
func (t *Tracer) ClientInterceptor() grpc.UnaryClientInterceptor {
	return otgrpc.OpenTracingClientInterceptor(t.Client)
}

// ServerInterceptor 會回傳一個可供 gRPC 中 NewServer 函式使用的 Interceptor。
func (t *Tracer) ServerInterceptor() grpc.UnaryServerInterceptor {
	return otgrpc.OpenTracingServerInterceptor(t.Client)
}

// StartSpan 會建立一個新的追蹤週期。
func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.Client.StartSpan(operationName, opts...)
}

// StartChildSpan 會接收一個父追蹤週期，並且建立一個子追蹤週期。
func (t *Tracer) StartChildSpan(operationName string, parent opentracing.Span) opentracing.Span {
	return t.Client.StartSpan(operationName, opentracing.ChildOf(parent.Context()))
}

// Inject 會在目前的追蹤週期中插入新的內容以利於傳遞到其他地方。
func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return t.Client.Inject(sm, format, carrier)
}

// Extract 會解開一個追蹤週期，並回傳其資料。
func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return t.Client.Extract(format, carrier)
}
