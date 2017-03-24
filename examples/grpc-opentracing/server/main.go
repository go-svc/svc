package main

import (
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/go-svc/svc/examples/grpc-fixed-lb/pb"
	"github.com/go-svc/svc/sd/lb"
	"github.com/go-svc/svc/tracing/opentracing"
)

// server 建構體會實作 Todo 的 gRPC 伺服器。
type server struct {
	db pb.TodoClient
}

// Add 會呼叫遠端 gRPC 資料庫並插入新的工作記事。
func (s *server) Add(ctx context.Context, in *pb.Task) (*pb.Task, error) {
	// 將接收到的資料透過 gRPC 客戶端傳送到遠端資料庫伺服器。
	//
	// 注意：這裡將接收到的 ctx 再次傳遞出去，所以追蹤器就能夠追蹤函式的動向。
	s.db.Add(ctx, in)
	return in, nil
}

// List 會取得遠端 gRPC 資料庫的所有工作記事並回傳至本地客戶端。
func (s *server) List(ctx context.Context, in *pb.Void) (*pb.Tasks, error) {
	// 透過 gRPC 客戶端呼叫遠端資料庫伺服器的 List 函式，
	// 用以取得工作記事列表。
	//
	// 注意：這裡將接收到的 ctx 再次傳遞出去，所以追蹤器就能夠追蹤函式的動向。
	tasks, _ := s.db.List(ctx, in)
	return tasks, nil
}

// newDB 會建立並回傳一個基於負載平衡的新 gRPC 客戶端連線到指定的 gRPC 資料庫伺服器。
func newDB(tracer *opentracing.Tracer) pb.TodoClient {
	// 取得實例群的地址，如果以 `go run` 執行，那麼實例群的地址就是 `os.Args[1]` 而不是 `os.Args[0]`。
	addrs := os.Args[0]
	if len(os.Args) == 2 {
		addrs = os.Args[1]
	}
	// 將實例群地址以逗點切分成字串陣列。
	instances := strings.Split(addrs, ", ")
	// 建立一個基於固定實例群的負載平衡器。
	balancer := lb.NewBalancer(lb.FixedOption{
		Mode:      lb.RoundRobin,
		Instances: instances,
	})

	// 透過負載平衡連線到資料庫，並且掛載追縱器。
	conn, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(balancer), grpc.WithUnaryInterceptor(tracer.ClientInterceptor()))
	if err != nil {
		log.Fatalf("連線失敗：%v", err)
	}

	return pb.NewTodoClient(conn)
}

// newTracer 會建立並回傳一個新的 Tracer，以用來定義此服務的追蹤器。
func newTracer(serviceName string) *opentracing.Tracer {
	tracer, err := opentracing.NewTracer(opentracing.Tracer{
		Collector: opentracing.Collector{
			URL: "http://localhost:9411/api/v1/spans",
		},
		Recorder: opentracing.Recorder{
			Debug:       false,
			Host:        "127.0.0.1:50051",
			ServiceName: serviceName,
		},
	})
	if err != nil {
		log.Fatalf("建立 OpenTracing 追蹤器時發生錯誤：%v", err)
	}

	return tracer
}

func main() {
	// 監聽指定埠口，這樣服務才能在該埠口執行。
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("無法監聽該埠口：%v", err)
	}

	// 建立一個此服務的追蹤器。
	tracer := newTracer("中央伺服器")
	// 建立另一個專門給資料庫客戶端的追蹤器。
	dbTracer := newTracer("資料庫客戶端")

	// 建立新 gRPC 伺服器並註冊 Todo 服務，並且帶入本服務的追蹤器，
	// 如此一來才可以追蹤此服務的所有動向。
	s := grpc.NewServer(grpc.UnaryInterceptor(tracer.ServerInterceptor()))
	pb.RegisterTodoServer(s, &server{
		// 建立連線到資料庫伺服器，所以稍後才能在本地伺服器中呼叫和資料庫相關的功能。
		// 並且傳入一個專門給資料庫客戶端用的追蹤器，所以我們才能夠掌握和資料庫的連線動向。
		db: newDB(dbTracer),
	})

	// 開始在指定埠口中服務。
	if err := s.Serve(lis); err != nil {
		log.Fatalf("無法提供服務：%v", err)
	}
}
