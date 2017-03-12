package lb

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-svc/svc/sd/consul"
	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

var (
	RoundRobin = 0
	Random     = 1
)

type ConsulOption struct {
	ServiceName string
	Mode        int
	Consul      consul.Client
	Tag         string
}

// ConsulBalancer 是一個負載平衡器。
type ConsulBalancer struct {
	Option      ConsulOption
	lastAddress int
	Addresses   []grpc.Address
}

// NewBalancer 會依照接收的選項型態來建立新的指定形態負載平衡器。
func NewBalancer(opt interface{}) grpc.Balancer {
	switch opt := opt.(type) {
	case ConsulOption:
		return &ConsulBalancer{Option: opt}
	}
	return nil
}

// Start 是用來初始化負載平衡器的函式。這個函式會在 gRPC 進行 Dial 時呼叫。
func (b ConsulBalancer) Start(target string, config grpc.BalancerConfig) error {

	return nil
}

func (b ConsulBalancer) Up(addr grpc.Address) (down func(error)) {

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
	if len(b.Addresses) == 0 {
		addr = grpc.Address{}
	} else {
		if b.Option.Mode == RoundRobin {
			if len(b.Addresses)-1 > b.lastAddress {
				addr = b.Addresses[b.lastAddress+1]
				b.lastAddress++
			} else {
				addr = b.Addresses[0]
				b.lastAddress = 0
			}
		} else {
			rand.Seed(time.Now().Unix())
			addr = b.Addresses[rand.Intn(len(b.Addresses))]
			fmt.Println(addr)
		}

	}
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
			services, _, _ := b.Option.Consul.Client().Catalog().Service(b.Option.ServiceName, b.Option.Tag, &api.QueryOptions{})
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
func (b ConsulBalancer) Close() error {

	return nil
}
