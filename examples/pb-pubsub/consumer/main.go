package main

import (
	"log"

	api "github.com/bitly/go-nsq"
	"github.com/go-svc/svc/examples/pb-pubsub/pb"
	"github.com/go-svc/svc/pubsub/nsq"
	"github.com/gogo/protobuf/proto"
)

func main() {
	// 建立「已接收」頻道，作為是否接收到訊息的一個開關。
	received := make(chan bool)
	// 建立新的 NSQ 客戶端。
	client := nsq.NewClient(nsq.Config{
		Producers: nsq.Producers{
			TCP:  "127.0.0.1:4150",
			HTTP: "127.0.0.1:4151",
		},
		Lookupds: []string{"127.0.0.1:4161"},
	}, api.NewConfig())

	// 建立一個基於 `msg` 主題的 `user` 頻道。
	ch := client.CreateChannel("msg", "user")
	// 以剛才建立的 `user` 頻道來建立一個消費者，用來「消化」基於 `msg` 主題的 `user` 頻道訊息。
	consumer, _ := client.NewConsumer(ch)

	// 建立訊息接收函式，當我們接收到訊息就會呼叫這個函式。
	consumer.AddHandler(func(msg *api.Message) error {
		var m pb.Msg
		// 將接收到的資料以 Protobuf 解碼。
		proto.Unmarshal(msg.Body, &m)
		// 顯示已經解碼的 Protobuf 內容。
		log.Printf("接收到了一個 Protobuf 資料，而 Content 是：%v", m.Content)

		// 對「已接收」頻道傳送 `true` 就能讓程式結束。
		received <- true
		return nil
	})

	// 連線到 NSQ 叢集，而不是單個 NSQ，這樣更安全與可靠。
	consumer.ConnectToNSQLookupds()
	// 除非接收到訊息，不然我們就讓程式卡住。
	<-received
}
