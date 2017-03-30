package main

import (
	"log"
	"net"

	"github.com/mwitkow/grpc-proxy/proxy"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// director 會負責指揮進入的代理請求應該流向到哪裡。
func director(ctx context.Context, fmn string) (*grpc.ClientConn, error) {
	// 將接下來的請求都流向到 `localhost:50051` 地址。
	return grpc.DialContext(ctx, "localhost:50051",
		grpc.WithCodec(proxy.Codec()),
		grpc.WithInsecure())
}

func main() {
	// 監聽指定埠口，這樣代理伺服器才能在該埠口執行。
	lis, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalln(err)
	}
	// 建立新的 gRPC 伺服器，如果接收到外來請求，
	// 而且這個請求是我們在這裡沒有定義的，
	// 那麼就呼叫 `director` 函式，然後決定這些請求應該流向到哪些服務。
	s := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))
	// 開始在指定埠口中服務。
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
