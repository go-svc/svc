package nsq

import (
	"fmt"
	"os/exec"

	api "github.com/bitly/go-nsq"
)

// Config 是 NSQ 建立客戶端時的設定檔。
type Config struct {
	// Producers 是連線到 NSQ 生產者的 HTTP、TCP 設定檔。
	Producers Producers
	// Lookupds 是用來設定單個或多個 NSQ 的 Lookupd 位址。
	Lookupds []string
}

// Client 是 NSQ 的客戶端。
type Client struct {
	Config    Config
	APIConfig *api.Config
	Producer  *api.Producer
}

// Producers 存放著訊息生產者的位址資訊。
type Producers struct {
	TCP  string
	HTTP string
}

// NewClient 會建立一個新的客戶端。
func NewClient(conf Config, apiConf *api.Config) *Client {
	prod, _ := api.NewProducer(conf.Producers.TCP, apiConf)

	return &Client{Config: conf, APIConfig: apiConf, Producer: prod}
}

// CreateTopic 會建立並返回一個話題。
func (c *Client) CreateTopic(name string) *Topic {
	cmd := exec.Command("curl", "-X", "POST", fmt.Sprintf("http://%s/topic/create?topic=%s", c.Config.Producers.HTTP, name))
	cmd.Start()
	cmd.Wait()

	return c.Topic(name)
}

// CreateChannel 會建立並返回一個基於指定話題的頻道。
func (c *Client) CreateChannel(topic, ch string) *Channel {
	cmd := exec.Command("curl", "-X", "POST", fmt.Sprintf("http://%s/channel/create?topic=%s&channel=%s", c.Config.Producers.HTTP, topic, ch))
	cmd.Start()
	cmd.Wait()

	return c.Channel(topic, ch)
}

// NewConsumer 會建立並返回一個基於指定頻道的消費者。
func (c *Client) NewConsumer(ch *Channel) (*Consumer, error) {
	consumer, err := api.NewConsumer(ch.Topic, ch.Name, c.APIConfig)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Topic:    ch.Topic,
		Channel:  ch.Name,
		Lookupds: c.Config.Lookupds,
		consumer: consumer,
	}, nil
}

// Channel 會回傳一個指定的頻道。
func (c *Client) Channel(topic, name string) *Channel {
	return &Channel{Topic: topic, Name: name, Producer: c.Producer}
}

// Topic 會回傳一個指定的話題。
func (c *Client) Topic(name string) *Topic {
	return &Topic{Name: name, Producer: c.Producer}
}

// Consumer 呈現了一個訊息消費者。
type Consumer struct {
	Topic    string
	Channel  string
	Lookupds []string
	Lookupd  string
	consumer *api.Consumer
}

// AddHandler 會替這個消費者新增一個處理函式，一但消費者接收到任何訊息就會呼叫這個函式。
func (c *Consumer) AddHandler(handler func(msg *api.Message) error) {
	c.consumer.AddHandler(api.HandlerFunc(func(m *api.Message) error {
		return handler(m)
	}))
}

// ConnectToNSQLookupds 會讓目前的消費者連線到 NSQLookupds 叢集。
func (c *Consumer) ConnectToNSQLookupds() {
	c.consumer.ConnectToNSQLookupds(c.Lookupds)
}

// Channel 呈現了一個基於特定話題的頻道。
type Channel struct {
	Topic    string
	Name     string
	Producer *api.Producer
}

// Topic 呈現了一個話題。
type Topic struct {
	Name     string
	Producer *api.Producer
}

// Publish 會在目前的話題上發布一個新的訊息。
func (t *Topic) Publish(body []byte) error {
	return t.Producer.Publish(t.Name, body)
}

// MultiPublish 可以在目前的話題上發佈多個訊息。
func (t *Topic) MultiPublish(body [][]byte) error {
	return t.Producer.MultiPublish(t.Name, body)
}
