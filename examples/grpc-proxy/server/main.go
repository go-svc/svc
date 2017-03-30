package main

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/context"

	"github.com/go-svc/svc/examples/grpc-proxy/pb"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

// server 建構體會實作 TestService 的 gRPC 伺服器。
type server struct{}

// Ping 會回傳接收到的訊息。
func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Value: in.Value}, nil
}

// PingStream 是一個串流，類似 WebSocket，會向客戶端每隔 200 毫秒就傳送 1 則訊息，總共會回傳 10 則訊息。
func (s *server) PingStream(ping *pb.PingRequest, stream pb.TestService_PingStreamServer) error {
	// 建立一個會跑 10 次的迴圈。
	for i := 0; i < 10; i++ {
		// 每隔 200 毫秒。
		<-time.After(200 * time.Millisecond)
		// 向串流發送訊息。
		stream.Send(&pb.PingResponse{
			Value:   ping.Value,
			Counter: int32(i),
		})
	}
	return nil
}

// director 會負責指揮進入的代理請求應該流向到哪裡。
func director(ctx context.Context, fmn string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, "localhost:50051",
		grpc.WithCodec(proxy.Codec()),
		grpc.WithInsecure())
}

func main() {
	// 監聽指定埠口，這樣服務才能在該埠口執行。
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("無法監聽該埠口：%v", err)
	}
	// 建立新 gRPC 伺服器並註冊 TestService 服務。
	s := grpc.NewServer(grpc.CustomCodec(proxy.Codec()))
	pb.RegisterTestServiceServer(s, &server{})
	// 透過 `proxy` 套件註冊一個服務，如此一來代理伺服器才能夠知道這個服務。
	proxy.RegisterService(s, director,
		"go-svc.svc.pb.TestService",
		"PingEmpty", "Ping", "PingError", "PingList", "PingStream",
	)
	// 開始在指定埠口中服務。
	if err := s.Serve(lis); err != nil {
		log.Fatalf("無法提供服務：%v", err)
	}
}
