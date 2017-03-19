# grpc-database

這個範例包含了資料庫和服務之間的溝通。這些服務皆共享同個 Protobuf 檔案，工作記事服務將會建立一個新的 gRPC 連線到資料庫服務，而客戶端將會呼叫工作記事服務。

## 範例

```bash
# 先啟動資料庫服務。
go run ./database/main.go
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
