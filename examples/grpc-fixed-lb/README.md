# grpc-fixed-lb

這個範例是基於 `grpc-database` 衍生出來的。這是個最簡單、基本的負載平衡範例，傳入固定的實例地址，就可以直接進行負載平衡。如果希望地址不要是固定的，可以參考：`grpc-lb`。

## 範例

```bash
# 啟動兩個資料庫服務實例，分別在 50050 和 50052 埠口上部署。
go run ./database-svc/server/main.go 50050
go run ./database-svc/server/main.go 50052
# 接著啟動工作記事服務並傳入兩個資料庫實例地址來進行負載平衡。
go run ./todo-svc/server/main.go "localhost:50050, localhost:50052"
# 然後開啟客戶端呼叫服務函式。
go run ./todo-svc/client/main.go
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
