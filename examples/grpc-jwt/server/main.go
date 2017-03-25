package main

import (
	"log"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/go-svc/svc/auth/jwt"
	"github.com/go-svc/svc/examples/grpc-metadata/pb"
	"github.com/go-svc/svc/metadata"
)

// server 建構體會實作 Calculator 的 gRPC 伺服器。
type server struct{}

// tokenContext 是 JSON Web Token 的內容。
type tokenContext struct {
	Iat      int64
	Nbf      int64
	Username string
}

// Plus 會將傳入的數字加總。
func (s *server) Plus(ctx context.Context, in *pb.CalcRequest) (*pb.CalcReply, error) {
	// 從 Context 裡面取得接收到的 JSON Web Token。
	md, _ := metadata.FromContext(ctx)
	token := md["jwt"][0]

	// 以指定的密碼建立新的 JSON Web Token 解析器。
	jwtParser := jwt.NewParser("JWT-Secret")
	// 解析 JSON Web Token。
	var tokenCtx tokenContext
	err := jwtParser.Parse(token, &tokenCtx)
	if err != nil {
		log.Fatalf("解析 JWT 時發生錯誤：%v", err)
	}
	log.Printf("已解析 JSON Web Token，其中的 Username 是：%v\n", tokenCtx.Username)

	// 計算傳入的數字。
	result := in.NumberA + in.NumberB
	// 包裝成 Protobuf 建構體並回傳。
	return &pb.CalcReply{Result: result}, nil
}

func main() {
	// 監聽指定埠口，這樣服務才能在該埠口執行。
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("無法監聽該埠口：%v", err)
	}

	// 建立新 gRPC 伺服器並註冊 Calculator 服務。
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	// 開始在指定埠口中服務。
	if err := s.Serve(lis); err != nil {
		log.Fatalf("無法提供服務：%v", err)
	}
}
