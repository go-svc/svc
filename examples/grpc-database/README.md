# grpc-database

這個範例包含了資料庫和服務之間的溝通。這些服務皆共享同個 Protobuf 檔案，工作記事服務將會建立一個新的 gRPC 連線到資料庫服務，而客戶端將會呼叫工作記事服務。

## 範例

```bash
# 先啟動資料庫服務。
go run ./database-svc/server/main.go
# 接著是工作記事服務。
go run ./todo-svc/server/main.go
# 然後開啟客戶端呼叫服務函式。
go run ./todo-svc/client/main.go
```

```bash
2017/03/12 06:35:38 Add 結果：測試, 這是測試的內容。
2017/03/12 06:35:38 List 結果：{"tasks":[{"title":"測試","content":"這是測試的內容。"}]}
```