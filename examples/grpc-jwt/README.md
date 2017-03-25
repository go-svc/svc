# grpc-jwt

這個 gRPC 範例衍生自 `grpc-metadata`，這會以一個特定的密碼簽署 [JSON Web Token](http://jwt.io/) 資料，並夾帶在中繼資料中傳遞至伺服端，伺服端接收到該資料後以相同的密碼解析該 JSON Web Token 以檢驗是否可用。

## 範例

```bash
# 啟動 gRPC 伺服器。
go run ./server/main.go
# 開啟客戶端與伺服器溝通。
go run ./client/main.go
```

```bash
2017/03/25 18:47:52 已簽發 JSON Web Token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTA0Mzg4NzIsIm5iZiI6MTQ5MDQzODg3MiwidXNlcm5hbWUiOiJZYW1pT2R5bWVsIn0.eME34UQ5s4QeqCQk_lgLhvkrELxyqzsmr_yIWdYWjXg
2017/03/25 18:47:52 回傳結果：64

```

```bash
2017/03/25 18:47:52 已解析 JSON Web Token，其中的 Username 是：YamiOdymel
```

## Protobuf

```proto
// Calculator 定義了一個計算用的服務。
service Calculator {
    rpc Plus (CalcRequest) returns (CalcReply) {}
}

// CalcRequest 包含了兩個數字，將會傳送至計算服務並對兩個數字進行計算。
message CalcRequest {
    int32 number_a = 1;
    int32 number_b = 2;
}

// CalcReply 是計算結果，將會回傳給客戶端。
message CalcReply {
    int32 result = 1;
}
```
