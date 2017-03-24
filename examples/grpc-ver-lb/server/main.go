package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/go-svc/svc/examples/grpc-ver-lb/pb"
	"github.com/go-svc/svc/sd/consul"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
)

// port 是指定部署的埠口。
var port string

// version 是這個服務的版本號。
var version string

// server 建構體會實作 Test 的 gRPC 伺服器。
type server struct {
}

// Ping 會返回一個目前服務版本號的內容。
func (s *server) Ping(ctx context.Context, in *pb.Void) (*pb.Info, error) {
	return &pb.Info{Version: version}, nil
}

// registerService 會將此服務註冊到服務探索中心。
func registerService() {
	// 以預設設定檔建立 Consul 客戶端。
	sd, _ := consul.NewClient(api.DefaultConfig())
	// 取得部署的埠口。
	p, _ := strconv.Atoi(port)
	// 註冊資料庫服務到服務探索中心。
	sd.Register(&api.AgentServiceRegistration{
		ID:   uuid.NewV4().String(),
		Name: "Test",
		Port: p,
		// 將版本號帶入標籤中，所以我們就能在客戶端透過標籤篩選服務的版本。
		Tags: []string{"test", version},
	})
}

func main() {
	// 取得部署的埠口，如果以 `go run` 執行，那麼埠口就是 `os.Args[1]` 而不是 `os.Args[0]`。
	port = os.Args[0]
	if len(os.Args) == 3 {
		port = os.Args[1]
	}
	// 取得部署的版本號。
	version = os.Args[1]
	if len(os.Args) == 3 {
		version = os.Args[2]
	}

	// 監聽指定埠口，這樣服務才能在該埠口執行。
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("無法監聽該埠口：%v", err)
	}

	// 建立新 gRPC 伺服器並註冊 Test 服務。
	s := grpc.NewServer()
	pb.RegisterTestServer(s, &server{})

	// 將此服務註冊到服務探索中心。
	registerService()

	// 開始在指定埠口中服務。
	if err := s.Serve(lis); err != nil {
		log.Fatalf("無法提供服務：%v", err)
	}
}
