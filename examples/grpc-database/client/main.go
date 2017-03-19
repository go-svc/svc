package main

import (
	"encoding/json"
	"log"

	"github.com/go-svc/svc/examples/grpc-database/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// 連線到遠端 gRPC 伺服器。
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
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
