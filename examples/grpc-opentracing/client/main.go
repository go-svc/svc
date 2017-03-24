package main

import (
	"encoding/json"
	"log"

	"github.com/go-svc/svc/examples/grpc-fixed-lb/pb"
	"github.com/go-svc/svc/tracing/opentracing"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// 建立新的 Tracer，也就是追蹤器。
	tracer, err := opentracing.NewTracer(opentracing.Tracer{
		Collector: opentracing.Collector{
			URL: "http://localhost:9411/api/v1/spans",
		},
		Recorder: opentracing.Recorder{
			Debug:       false,
			Host:        "127.0.0.1:0",
			ServiceName: "客戶端",
		},
	})
	if err != nil {
		log.Fatalf("建立 OpenTracing 追蹤器時發生錯誤：%v", err)
	}

	// 連線到遠端 gRPC 伺服器，
	// 並傳入 Tracer，如此一來就能夠追蹤程式接下來的行蹤。
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithUnaryInterceptor(tracer.ClientInterceptor()))
	if err != nil {
		log.Fatalf("連線失敗：%v", err)
	}
	defer conn.Close()

	// 建立新的 Todo 客戶端，所以等一下就能夠使用 Todo 的所有方法。
	c := pb.NewTodoClient(conn)
	// 傳送新請求到遠端 gRPC 伺服器 Todo 中，並呼叫 Add 函式，新增一個工作記事。
	r, err := c.Add(context.Background(), &pb.Task{Title: "測試", Content: "這是測試的內容。"})
	if err != nil {
		log.Fatalf("無法執行 Add 函式：%v", err)
	}
	log.Printf("Add 結果：%v, %v", r.Title, r.Content)

	// 傳送新請求到遠端 gRPC 伺服器 Todo 中，並呼叫 List 函式，列出所有在資料庫內的工作記事。
	tasks, err := c.List(context.Background(), &pb.Void{})
	if err != nil {
		log.Fatalf("無法執行 List 函式：%v", err)
	}
	t, _ := json.Marshal(tasks)
	log.Printf("List 結果：%v", string(t))
}
