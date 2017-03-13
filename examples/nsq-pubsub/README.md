# nsq-pubsub

通常在撰寫微服務時，我們希望服務之間可以是異步進行的，所以我們會透過像是 [NSQ](http://nsq.io/)、[RabbitMQ](http://www.rabbitmq.com/) 這樣的「分布式訊息佇列系統」來達成這項功能。這個理論非常簡單，一個服務發送訊息，另一個負責接收。因為不用在意另一個服務是不是已經完成了，所以我們就可以直接進行下一步動作來增加處理速度與效率。

這個範例需要安裝 [NSQ](http://nsq.io/)。

## 範例

```bash
# 啟動生產者，這會向 NSQ 傳送一個訊息到 `msg` 主題。
go run ./producer/main.go
# 啟動消費者，這會向 NSQ 註冊一個基於 `msg` 主題的 `user` 頻道，
# 所以就能接收 `msg` 主題的訊息。
go run ./consumer/main.go
```

```bash
2017/03/14 04:17:09 INF    2 [msg/user] querying nsqlookupd http://127.0.0.1:4161/lookup?topic=msg
2017/03/14 04:17:09 INF    2 [msg/user] (127.0.0.1:4150) connecting to nsqd
2017/03/14 04:17:14 接收到了一個訊息：Hello, world!
```