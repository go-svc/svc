# es-pubsub

這個範例詮釋了如何與 [EventStore](https://geteventstore.com/)（事件存儲中心）進行互動，事件存儲中心就像是 NSQ 那樣的分布式訊息佇列系統，但差別是事件存儲中心會將傳遞的永久事件保存下來，並且能在服務重新上線時重播所有的事件，所以服務也就能夠擁有以往所有的資料。

這個範例需要安裝 [EventStore](https://geteventstore.com/)。

## 範例

```bash
# 啟動串流撰寫者，這會向 EventStore 的 `user.created` 串流發送事件。
go run ./writer/main.go
# 啟動接收者，這會不斷地監聽 EventStore 中的 `user.created` 串流，
# 並且處理接收到的事件。
go run ./reader/main.go
```

```bash
事件已接收：小安, 再單身超過十幾年，小安就能夠成為魔法師了。
```