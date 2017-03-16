package main

import (
	"log"

	"github.com/go-svc/svc/examples/grpc-ver-lb/pb"
	"github.com/go-svc/svc/sd/consul"
	"github.com/go-svc/svc/sd/lb"
	"github.com/go-svc/svc/version"
	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// 建立新的 Consul 客戶端。
	sd, _ := consul.NewClient(api.DefaultConfig())
	// 建立一個基於 Consul 的負載平衡器。
	balancer := lb.NewBalancer(lb.ConsulOption{
		Name: "Test",
		Mode: lb.RoundRobin,
		// 指定期望的服務版本。
		Tag:    version.Define(1, 0, 0, "stable").String(),
		Client: sd,
	})

	// 連線到遠端 gRPC 伺服器並採用負載平衡器。
	conn, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(balancer))
	if err != nil {
		log.Fatalf("連線失敗：%v", err)
	}
	defer conn.Close()

	// 建立新的 Test 客戶端，所以等一下就能夠使用 Test 的所有方法。
	c := pb.NewTestClient(conn)
	// 傳送新請求到遠端 gRPC 伺服器 Test 中，並呼叫 Ping 函式。
	r, err := c.Ping(context.Background(), &pb.Void{})
	if err != nil {
		log.Fatalf("無法執行 Ping 函式：%v", err)
	}
	log.Printf("Ping 結果：%s", r.Version)
}
