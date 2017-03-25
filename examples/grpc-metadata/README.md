# grpc-metadata

這是基本的 gRPC 範例加上傳遞、解析中繼資料，這能夠讓你在傳遞資料時附上一些額外的資料（如：客戶端 IP、名稱⋯⋯等）。

## 範例

```bash
# 啟動 gRPC 伺服器。
go run ./server/main.go
# 開啟客戶端與伺服器溝通。
go run ./client/main.go
```

```bash
2017/03/25 18:47:52 回傳結果：64
```

```bash
2017/03/25 18:47:52 已接收到中繼資料，資料來源是：CalcClient
2017/03/25 18:47:52 中繼資料 `test_meta` 的內容是：這是個測試用的中繼資料。
2017/03/25 18:47:52 中繼資料 `test_meta2` 的內容是：你能夠透過中繼資料
2017/03/25 18:47:52 中繼資料 `test_meta3` 的內容是：來在資料傳遞時夾帶一些額外的有用資訊。
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
