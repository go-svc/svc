package lb

import (
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

var (
	// RoundRobin 是輪詢模式，負載平衡器每次都會按照順序逐一採用不同的伺服器，確保每個伺服器都會被使用到。
	RoundRobin = 0
	// Random 是隨機模式，負載平衡器每次都會隨機挑選一個伺服器使用。
	Random = 1
)

// NewBalancer 會依照接收的選項型態來建立新的指定形態負載平衡器。
func NewBalancer(opt interface{}) grpc.Balancer {
	switch opt := opt.(type) {
	// 接收到 ConsulOption 就建立一個 Consul 負載平衡器。
	case ConsulOption:
		return &ConsulBalancer{Option: opt}
	// 接收到 FixedOption 就建立一個固定的負載平衡器。
	case FixedOption:
		return &FixedBalancer{Option: opt}
	}
	return nil
}

// pick 會依照傳入的模式來在傳入的實例群中挑選出一個伺服器，並將其與新的最後索引一同回傳。
func pick(mode int, addrs []grpc.Address, last int) (grpc.Address, int) {
	// 如果實例群是空的，就回傳空的地址。
	if len(addrs) == 0 {
		return grpc.Address{}, 0
	}

	// 依照指定的模式來在實力群內進行挑選。
	switch mode {
	// 輪詢模式。
	case RoundRobin:
		// 如果實例群中還有下一個可供使用，那麼就選擇下一個。
		if len(addrs)-1 > last {
			return addrs[last+1], last + 1
		}
		// 否則從頭來過。
		return addrs[0], 0
	// 隨機模式。
	default:
		rand.Seed(time.Now().Unix())
		return addrs[rand.Intn(len(addrs))], 0
	}
}
