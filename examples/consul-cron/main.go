package main

import (
	"log"

	"github.com/go-svc/svc/cron"
)

func main() {
	// 建立新的 Consul 客戶端。
	//sd, _ := consul.NewClient(api.DefaultConfig())

	//cr := cron.NewSchedulerWithStore(cron.NewConsul(sd))
	//cr.Every(1)

	c := cron.NewScheduler()
	c.Job("sec").Every(1).Second().Do(func() {
		log.Println("排程：每 1 秒。")
	})
	c.Job("3-sec").Every(3).Seconds().Do(func() {
		log.Println("排程：每 3 秒，移除每 1 秒排程。")
		c.Remove("sec")
	})
	c.Job("min").Every(1).Minute().Do(func() {
		log.Println("排程：每 1 分鐘。")
	})
	log.Println("排程開始。")
	c.Start()
	log.Println("優先執行 `min` 排程，而不是等到過了 1 分鐘才執行。")
	c.Run("min")
}
