package main

import (
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/go-svc/svc/examples/grpc-database/pb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"google.golang.org/grpc"
)

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
	// 監聽指定埠口，這樣服務才能在該埠口執行。
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("無法監聽該埠口：%v", err)
	}

	// 建立新 gRPC 伺服器並註冊 Todo 服務。
	s := grpc.NewServer()
	pb.RegisterTodoServer(s, &server{
		// 建立連線到資料庫伺服器，所以稍後才能在本地伺服器中呼叫和資料庫相關的功能。
		db: newConnection(),
	})

	// 開始在指定埠口中服務。
	if err := s.Serve(lis); err != nil {
		log.Fatalf("無法提供服務：%v", err)
	}
}
