package main

import (
	api "github.com/bitly/go-nsq"
	"github.com/go-svc/svc/examples/pb-pubsub/pb"
	"github.com/go-svc/svc/pubsub/nsq"
	"github.com/gogo/protobuf/proto"
)

func main() {
	// 建立新的 NSQ 客戶端。
	client := nsq.NewClient(nsq.Config{
		Producers: nsq.Producers{
			TCP:  "127.0.0.1:4150",
			HTTP: "127.0.0.1:4151",
		},
		Lookupds: []string{"127.0.0.1:4161"},
	}, api.NewConfig())

	// 建立一個新的 `msg` 主題。
	topic := client.CreateTopic("msg")

	// 建立一個名為 Msg 的 Protobuf 資料。
	data := pb.Msg{
		Content: "你好，世界！",
	}
	// 將資料編碼成 Protocol Buffer 格式（請注意是傳入 Pointer）。
	dataBuffer, _ := proto.Marshal(&data)

	// 在 `msg` 主題中傳入已經編碼的 Protobuf 內容。
	topic.Publish(dataBuffer)
}
