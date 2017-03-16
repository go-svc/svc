# grpc-ver-lb

這個範例是基於 `grpc-lb` 衍生出來的。但這個範例的服務探索支援版本控制，這意味著你可以控制服務僅能呼叫特定版本的服務用以避免升級時相容性問題。在這個範例中雖然有兩個服務，但客戶端僅會呼叫 `1.0.0+stable` 版本的服務，因為我們已經指定了。

## 範例

```bash
# 啟動兩個伺服器於 50050 和 50052 埠口，但設定成不同版本。
go run ./server/main.go 50050 "1.0.0+stable"
go run ./server/main.go 50052 "2.0.0+stable"
# 然後開啟客戶端呼叫服務函式。
go run ./client/main.go
```

```bash
2017/03/17 01:53:46 Ping 結果：1.0.0+stable
```

## Protobuf

```proto
// Test 是一個用來測試是否有回應的服務。
service Test {
    rpc Ping(Void)  returns (Info) {}
}

// Void 呈現一個什麼都沒有的資料。
message Void {
}

// Info 呈現一個帶有服務版本號的資料。
message Info {
    string version = 1;
}
```
