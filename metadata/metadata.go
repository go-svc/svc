package metadata

import (
	"golang.org/x/net/context"
	md "google.golang.org/grpc/metadata"
)

// Info 是用來建立中繼資料的資訊建構體。
type Info struct {
	ServiceName string
}

// Metadata 呈現了一個基本的中繼資料。
type Metadata struct {
	// MD 存放了原生的 gRPC metadata。
	MD *MD
}

// MD 是中繼資料的型態。
type MD map[string][]string

// New 會接收一個資訊建構體用來初始化、回傳一個新的中繼資料。
func New(info Info) *Metadata {
	m := md.New(map[string]string{
		"service_name": info.ServiceName,
	})
	// 把 md.MD 轉換成本地 MD 型態。
	metadata := MD(m)

	return &Metadata{
		MD: &metadata,
	}
}

// NewMap 會接收一個基於 map 結構的中繼資料然後轉換成 gRPC 的中繼資料。
func NewMap(metadata map[string]string) md.MD {
	return md.New(metadata)
}

// FromContext 會解析傳入的 Context 並取得、回傳其中的中繼資料。
func FromContext(ctx context.Context) (metadata MD, ok bool) {
	m, o := md.FromContext(ctx)
	metadata, ok = MD(m), o
	return
}

// ContextMap 會接收 map 型態的中繼資料並與目前的中繼資料合併，
// 然後將其轉換成 context.Context 以便利於傳入其他函式中。
func (m *Metadata) ContextMap(metadata map[string]string) context.Context {
	return md.NewContext(context.Background(), md.Join(md.MD(*m.MD), NewMap(metadata)))
}

// Context 會將目前的中繼資料轉換成 context.Context 以便利於傳入其他函式中。
func (m *Metadata) Context() context.Context {
	return md.NewContext(context.Background(), md.MD(*m.MD))
}
