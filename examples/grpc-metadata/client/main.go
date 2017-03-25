package main

import (
	"log"

	"github.com/go-svc/svc/examples/grpc-metadata/pb"
	"github.com/go-svc/svc/metadata"
	"google.golang.org/grpc"
)

func main() {
	// 連線到遠端 gRPC 伺服器。
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("連線失敗：%v", err)
	}
	defer conn.Close()

	// 建立一個新的中繼資料，並帶有客戶端服務名稱以方便伺服端辨識資料來源。
	md := metadata.New(metadata.Info{
		ServiceName: "CalcClient",
	})
	// 在中繼資料中插入自訂資料，稍後會在伺服端進行解析。
	ctx := md.ContextMap(map[string]string{
		"test_meta":  "這是個測試用的中繼資料。",
		"test_meta2": "你能夠透過中繼資料",
		"test_meta3": "來在資料傳遞時夾帶一些額外的有用資訊。",
	})

	// 建立新的 Calculator 客戶端，所以等一下就能夠使用 Calculator 的所有方法。
	c := pb.NewCalculatorClient(conn)
	// 傳送新請求到遠端 gRPC 伺服器 Calculator 中，呼叫 Plus 函式並夾帶中繼資料，讓兩個數字相加。
	r, err := c.Plus(ctx, &pb.CalcRequest{NumberA: 32, NumberB: 32})
	if err != nil {
		log.Fatalf("無法執行 Plus 函式：%v", err)
	}
	log.Printf("回傳結果：%d", r.Result)
}
