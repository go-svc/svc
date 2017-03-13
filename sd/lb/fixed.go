package lb

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// FixedOption 是用來建立基本、固定實例的負載平衡器選項。
type FixedOption struct {
	// Mode 是查詢模式，如 `lb.RoundRobin` 或 `lb.Random`。
	Mode int
	// Instances 集合了所有實例的地址，帶有位置和埠口。例如：`localhost:50050`。
	Instances []string
}

// FixedBalancer 是一個基本並帶有固定實例的負載平衡器。
type FixedBalancer struct {
	// Option 是負載平衡器建立時的選項。
	Option FixedOption
	// lastAddress 是上一次輪詢的地址索引，以此來確定下一個地址為何。
	lastAddress int
	// Addresses 包含了所有可用的實例地址。
	Addresses []grpc.Address
}

// Start 是用來初始化負載平衡器的函式。這個函式會在 gRPC 進行 Dial 時呼叫。
func (b *FixedBalancer) Start(target string, config grpc.BalancerConfig) error {

	return nil
}

func (b *FixedBalancer) Up(addr grpc.Address) (down func(error)) {
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
func (b *FixedBalancer) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	addr, last := pick(b.Option.Mode, b.Addresses, b.lastAddress)
	b.lastAddress = last

	put = func() {

	}
	err = nil
	return
}

func (b *FixedBalancer) Notify() <-chan []grpc.Address {
	ch := make(chan []grpc.Address)

	go func() {

		var address []grpc.Address
		for _, addr := range b.Option.Instances {
			address = append(address, grpc.Address{
				Addr: addr,
			})
		}

		b.Addresses = address

		ch <- b.Addresses
	}()

	return ch
}

// Close 會關閉負載平衡器。
func (b *FixedBalancer) Close() error {

	return nil
}
