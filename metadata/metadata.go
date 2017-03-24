package metadata

import (
	"golang.org/x/net/context"
	md "google.golang.org/grpc/metadata"
)

type Info struct {
	ServiceName string
}

type Metadata struct {
	MD *md.MD
}

type MD map[string][]string

func New(info Info) *Metadata {
	metadata := md.New(map[string]string{
		"service_name": info.ServiceName,
	})

	return &Metadata{
		MD: &metadata,
	}
}

func NewMap(metadata map[string]string) md.MD {
	return md.New(metadata)
}

func FromContext(ctx context.Context) (MD, bool) {
	metadata, ok := md.FromContext(ctx)
	return MD(metadata), ok
}

func (m *Metadata) Context(metadata map[string]string) context.Context {
	return md.NewContext(context.Background(), md.Join(*m.MD, NewMap(metadata)))
}
