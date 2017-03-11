# grpc

這是基本的 gRPC 範例。

## 範例

```bash
# 啟動 gRPC 伺服器。
go run ./server/main.go
# 開啟客戶端與伺服器溝通。
go run ./client/main.go
```

```bash
2017/03/12 06:37:53 回傳結果：64
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
