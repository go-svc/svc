package main

import "github.com/go-svc/svc/pubsub/eventstore"

func main() {
	// 建立新的 EventStore 客戶端。
	client := eventstore.NewClient(eventstore.Config{
		URL:      "http://127.0.0.1:2113",
		Username: "admin",
		Password: "changeit",
	})
	// 建立新的 `user.created` 串流。
	stream := client.CreateStream("user.created")
	// 在該串流上推送新的事件。
	stream.Append(eventstore.Event{
		Data: map[string]string{
			"nickname": "小安",
			"about_me": "再單身超過十幾年，小安就能夠成為魔法師了。",
		},
	})
}
