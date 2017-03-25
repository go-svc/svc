package main

import (
	"log"
	"time"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-svc/svc/auth/jwt"
	"github.com/go-svc/svc/examples/grpc-metadata/pb"
	"github.com/go-svc/svc/metadata"
	"google.golang.org/grpc"
)

// tokenContext 是 JSON Web Token 的內容。
type tokenContext struct {
	Iat      int64
	Nbf      int64
	Username string
}

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

	// 建立新的簽署者，其密碼是 `JWT-Secret`，並採用 HS256 演算法。
	signer := jwt.NewSigner("JWT-Secret", stdjwt.SigningMethodHS256)
	// 以剛才建立的簽署者簽署下列資料。
	token, err := signer.SignWithStruct(tokenContext{
		Nbf:      time.Now().Unix(),
		Iat:      time.Now().Unix(),
		Username: "YamiOdymel",
	})
	if err != nil {
		log.Fatalf("簽署 JSON Web Token 時發生錯誤：%v", err)
	}
	log.Printf("已簽發 JSON Web Token：%v", token)

	// 在中繼資料中插入名為 `jwt` 的 JSON Web Token，稍後會在伺服端進行解析。
	ctx := md.ContextMap(map[string]string{
		"jwt": token,
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
