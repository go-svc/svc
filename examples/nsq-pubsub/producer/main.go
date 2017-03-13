package main

import (
	api "github.com/bitly/go-nsq"
	"github.com/go-svc/svc/pubsub/nsq"
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
	// 在 `msg` 主題中發表 `Hello, world!` 內容。
	topic.Publish([]byte("Hello, world!"))
}
