# grpc-lb

這個範例是基於 `grpc-database` 衍生出來的。這個範例會需要安裝 [Consul](https://www.consul.io/)（如果不想安裝 Consul，可查看基本的負載平衡範例：`grpc-fixed-lb`），我們會透過 Consul 在進行負載平衡，資料庫皆會自動向服務探索中心註冊，而工作記事服務需要呼叫資料庫服務時，就會自動透過負載平衡器取得一個可用的資料庫服務實例。負載平衡器支援輪詢（Round-Robin）和隨機（Random）模式。

## 範例

```bash
# 啟動兩個資料庫服務實例，分別在 50050 和 50052 埠口上部署。
go run ./database/main.go 50050
go run ./database/main.go 50052
# 接著是工作記事服務。
go run ./server/main.go
# 然後開啟客戶端呼叫服務函式。
go run ./client/main.go
```

```bash
2017/03/12 06:35:38 Add 結果：測試, 這是測試的內容。
2017/03/12 06:35:38 List 結果：{"tasks":[{"title":"測試","content":"這是測試的內容。"}]}
```

## Protobuf

```proto
// Todo 是一個提供存放工作記事的服務。
service Todo {
    rpc Add(Task)  returns (Task)  {}
    rpc List(Void) returns (Tasks) {}
}

// Void 呈現一個什麼都沒有的資料。
message Void {
}

// Task 是單個工作記事資料。
message Task {
    string title   = 1;
    string content = 2;
}

// Tasks 會回傳多個工作記事資料。
message Tasks {
    repeated Task tasks = 1;
}
```
