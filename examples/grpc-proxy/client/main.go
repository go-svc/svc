package main

import (
	"io"
	"log"

	"github.com/go-svc/svc/examples/grpc-proxy/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// 連線到遠端 gRPC 「代理」伺服器。
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("無法連線到 Proxy：%v", err)
	}
	defer conn.Close()

	// 建立新的 TestService 客戶端，所以等一下就能夠使用 TestService 的所有方法。
	c := pb.NewTestServiceClient(conn)
	// 傳送新請求到遠端 gRPC 伺服器 TestService 中，並呼叫 Ping 函式。
	r, err := c.Ping(context.Background(), &pb.PingRequest{Value: "Hello, World!"})
	if err != nil {
		log.Fatalf("無法執行 Ping 函式：%v", err)
	}
	log.Printf("回傳結果：%s", r.Value)

	// 傳送新請求到遠端 gRPC 伺服器 TestService 中，並呼叫 PingStream 函式開始一個新的串流。
	stream, err := c.PingStream(context.Background(), &pb.PingRequest{Value: "Foobar!"})
	if err != nil {
		log.Fatalf("無法執行 PingList 函式：%v", err)
	}
	// 不斷地讀取串流。
	for {
		resp, err := stream.Recv()
		// io.EOF 表示這是串流的結尾，所以我們應該就這樣結束這個迴圈。
		if err == io.EOF {
			log.Println("PingList 串流已經到結尾，自動結束。")
			break
		} else if err != nil {
			log.Fatalf("讀取 PingList 串流時發生錯誤：%v", err)
		}
		log.Printf("回傳結果：%v", resp)
	}
}
