package main

import (
	"log"
	"net"
	"os"

	"golang.org/x/net/context"

	"github.com/go-svc/svc/examples/grpc-fixed-lb/pb"
	"github.com/go-svc/svc/tracing/opentracing"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"google.golang.org/grpc"
)

// port 是指定部署的埠口。
var port string

// server 建構體會實作 Todo 的 gRPC 伺服器。
type server struct {
	db *gorm.DB
}

// Add 會插入接收到的資料到本地資料庫。
func (s *server) Add(ctx context.Context, in *pb.Task) (*pb.Task, error) {
	// 在資料庫內插入剛才接收到的工作記事資料。
	s.db.Create(in)
	return in, nil
}

// List 會回傳本地資料庫中所有的工作記事。
func (s *server) List(ctx context.Context, in *pb.Void) (*pb.Tasks, error) {
	var tasks []*pb.Task
	// 取得資料庫內 `tasks` 資料表的所有紀錄。
	s.db.Find(&tasks)
	return &pb.Tasks{Tasks: tasks}, nil
}

// newConnection 會建立並回傳一個新的資料庫連線。
func newConnection() *gorm.DB {
	db, _ := gorm.Open("sqlite3", "/tmp/gorm.db")
	// 初始化資料表格
	db.DropTableIfExists(&pb.Task{})
	db.CreateTable(&pb.Task{})
	return db
}

func main() {
	// 取得部署的埠口，如果以 `go run` 執行，那麼埠口就是 `os.Args[1]` 而不是 `os.Args[0]`。
	port = os.Args[0]
	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	// 監聽指定埠口，這樣服務才能在該埠口執行。
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("無法監聽該埠口：%v", err)
	}

	// 建立一個給資料庫的追蹤器。
	tracer, err := opentracing.NewTracer(opentracing.Tracer{
		Collector: opentracing.Collector{
			URL: "http://localhost:9411/api/v1/spans",
		},
		Recorder: opentracing.Recorder{
			Debug:       false,
			Host:        "127.0.0.1:" + port,
			ServiceName: "資料庫",
		},
	})
	if err != nil {
		log.Fatalf("建立 OpenTracing 追蹤器時發生錯誤：%v", err)
	}

	// 建立新 gRPC 伺服器並註冊 Todo 服務。
	// 並且插入一個追蹤器以利於後續的動向追蹤。
	s := grpc.NewServer(grpc.UnaryInterceptor(tracer.ServerInterceptor()))
	pb.RegisterTodoServer(s, &server{
		// 建立連線到資料庫伺服器，所以稍後才能在本地伺服器中呼叫和資料庫相關的功能。
		db: newConnection(),
	})

	// 開始在指定埠口中服務。
	if err := s.Serve(lis); err != nil {
		log.Fatalf("無法提供服務：%v", err)
	}
}
