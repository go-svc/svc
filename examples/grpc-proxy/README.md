# grpc-proxy

這是 gRPC 和 Proxy 代理伺服器的範例，這讓你有一個代理伺服器用來隱藏內部、後端的伺服器，當客戶端要與伺服器有所溝通時都會透過 Proxy 伺服器。

## 範例

```bash
# 啟動 gRPC 伺服器，這個伺服器會處理所有的邏輯。
go run ./server/main.go
# 啟動 Proxy 伺服器，所有連入的連線都會轉交給 gRPC 伺服器。
go run ./proxy/main.go
# 開啟客戶端與 Proxy 伺服器溝通，這些訊息都會由 Proxy 轉交給 gRPC 伺服器。
go run ./client/main.go
```

```bash
2017/03/30 22:48:10 回傳結果：Hello, World!
2017/03/30 22:48:10 回傳結果：value:"Foobar!"
2017/03/30 22:48:10 回傳結果：value:"Foobar!" counter:1
2017/03/30 22:48:10 回傳結果：value:"Foobar!" counter:2
2017/03/30 22:48:10 回傳結果：value:"Foobar!" counter:3
2017/03/30 22:48:11 回傳結果：value:"Foobar!" counter:4
2017/03/30 22:48:11 回傳結果：value:"Foobar!" counter:5
2017/03/30 22:48:11 回傳結果：value:"Foobar!" counter:6
2017/03/30 22:48:11 回傳結果：value:"Foobar!" counter:7
2017/03/30 22:48:11 回傳結果：value:"Foobar!" counter:8
2017/03/30 22:48:12 回傳結果：value:"Foobar!" counter:9
2017/03/30 22:48:12 PingList 串流已經到結尾，自動結束。
```

## Protobuf

```proto
// PingRequest 是 Ping 時需要傳入的基本資料。
message PingRequest {
  string value = 1;
}

// PingResponse 呈現了一個 Ping 會回傳的訊息。
message PingResponse {
  string value   = 1;
  int32  counter = 2;
}

service TestService {
  // Ping 會接收 PingRequest 並將接收到的訊息轉換成 PingResponse 然後回傳。
  rpc Ping(PingRequest) returns (PingResponse) {}
  // PingStream 是一個串流，類似 WebSocket，會向客戶端每隔 200 毫秒就傳送 1 則 PingResponse，總共會回傳 10 次。
  rpc PingStream(PingRequest) returns (stream PingResponse) {}
}
```
