package main

import (
	"fmt"

	"github.com/go-svc/svc/pubsub/eventstore"
	"github.com/jetbasrawi/go.geteventstore"
)

func main() {
	// 建立「已接收」頻道，作為是否接收到訊息的一個開關。
	received := make(chan bool)
	// 建立新的 EventStore 客戶端。
	client := eventstore.NewClient(eventstore.Config{
		URL:      "http://127.0.0.1:2113",
		Username: "admin",
		Password: "changeit",
	})
	// 建立新的 `user.created` 串流。
	stream := client.CreateStream("user.created")

	// 在這個串流上建立相對應的處理函式，當接收到事件時就會呼叫這個函式。
	stream.AddHandler(func(r *goes.StreamReader) {
		// 建立空的事件內容和中繼資料變數，等一下會將接收到的事件解析到這兩個變數中。
		data := make(map[string]string)
		meta := make(map[string]string)
		// 呼叫 Scan 來將事件內容與中繼資料解析到變數中。
		r.Scan(&data, &meta)
		// 顯示解析後的資料。
		fmt.Printf("事件已接收：%s, %s", data["nickname"], data["about_me"])

		// 對「已接收」頻道傳送 `true` 就能讓程式結束。
		received <- true
	})

	// 除非接收到訊息，不然我們就讓程式卡住。
	<-received
}
