package lb

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/go-svc/svc/sd/consul"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

// ConsulOption 是用來建立基於 Consul 負載平衡器的選項。
type ConsulOption struct {
	// Name 是服務名稱，我們將會搜尋此服務的所有主機。
	Name string
	// Mode 是查詢模式，如 `lb.RoundRobin` 或 `lb.Random`。
	Mode int
	// Client 是 Go Svc 的 Consul 客戶端。
	Client consul.Client
	//
	Tag string
}

// ConsulBalancer 是一個基於 Consul 的負載平衡器。
type ConsulBalancer struct {
	// Option 是負載平衡器建立時的選項。
	Option ConsulOption
	// lastAddress 是上一次輪詢的地址索引，以此來確定下一個地址為何。
	lastAddress int
	// Addresses 包含了所有可用的實例地址。
	Addresses []grpc.Address
}

// Start 是用來初始化負載平衡器的函式。這個函式會在 gRPC 進行 Dial 時呼叫。
func (b *ConsulBalancer) Start(target string, config grpc.BalancerConfig) error {

	return nil
}

func (b *ConsulBalancer) Up(addr grpc.Address) (down func(error)) {

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
func (b *ConsulBalancer) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	addr, last := pick(b.Option.Mode, b.Addresses, b.lastAddress)
	b.lastAddress = last

	put = func() {

	}

	err = nil
	return
}

func (b *ConsulBalancer) Notify() <-chan []grpc.Address {
	ch := make(chan []grpc.Address)

	go func() {
		for {
			// 每隔一段時間就從服務探索中心裡取得可用的指定類別服務。
			<-time.After(3 * time.Second)

			// 從服務探索中心裡取得所有相關的服務。
			services, _, _ := b.Option.Client.Client().Catalog().Service(b.Option.Name, b.Option.Tag, &api.QueryOptions{})
			// 如果服務探索中心裡沒有需要的服務，則略過此次搜尋，直接進行下一輪搜尋。
			if len(services) == 0 {
				continue
			}

			var address []grpc.Address

			for _, svc := range services {
				address = append(address, grpc.Address{
					Addr: fmt.Sprintf("%s:%d", svc.Address, svc.ServicePort),
				})
			}

			b.Addresses = address

			ch <- b.Addresses
		}
	}()

	return ch
}

// Close 會關閉負載平衡器。
func (b *ConsulBalancer) Close() error {

	return nil
}
