package lb

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

type BalancerOption struct {
	ServiceName string
	Consul      *api.Client
	Tag         string
}

// Balancer 是一個負載平衡器。
type Balancer struct {
	Option    BalancerOption
	Addresses []grpc.Address
}

// NewBalancer 會建立新的負載平衡器。
func NewBalancer(opt BalancerOption) grpc.Balancer {

	return Balancer{Option: opt}
}

// Start 是用來初始化負載平衡器的函式。這個函式會在 gRPC 進行 Dial 時呼叫。
func (b Balancer) Start(target string, config grpc.BalancerConfig) error {

	return nil
}

func (b Balancer) Up(addr grpc.Address) (down func(error)) {

	down = func(err error) {

	}
	return
}

// Get 會基於 ctx 內容來回傳一個 gRPC 伺服器的地址。
//
// put 函式會在 RPC 完成或是失敗的時候呼叫。
// 這個函式可以用來收集該次 RPC 的狀態到負載平衡器做為統計。
//
// 當負載平衡器無法從錯誤中回復時應該回傳 err 錯誤，如果回傳了錯誤，gRPC 會將該次的 RPC 請求標記為失敗。
func (b Balancer) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	addr = grpc.Address{Addr: "localhost:50050"}

	put = func() {

	}
	err = nil
	return
}

func (b Balancer) Notify() <-chan []grpc.Address {
	ch := make(chan []grpc.Address)

	go func() {
		for {
			<-time.After(3 * time.Second)

			services, _, _ := b.Option.Consul.Catalog().Service(b.Option.ServiceName, b.Option.Tag, &api.QueryOptions{})

			var address []grpc.Address

			for _, svc := range services {
				address = append(address, grpc.Address{
					Addr: fmt.Sprintf("%s:%d", svc.Address, svc.ServicePort),
				})
			}
			//[]grpc.Address{
			//	{Addr: "localhost:50050"},
			//}
			b.Addresses = address

			ch <- b.Addresses
		}
	}()

	return ch
}

// Close 會關閉負載平衡器。
func (b Balancer) Close() error {

	return nil
}
