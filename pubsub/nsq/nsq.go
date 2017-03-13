package nsq

import (
	"fmt"
	"os/exec"

	api "github.com/bitly/go-nsq"
)

type Client interface {
	CreateTopic(name string) *Topic
	CreateChannel(topic, ch string) *Channel
	Topic(name string) *Topic
	Channel(topic, name string) *Channel
	NewConsumer(ch *Channel) (*Consumer, error)
}

type Config struct {
	Producers Producers
	Lookupds  []string
}

type client struct {
	Config
	APIConfig *api.Config
	Producer  *api.Producer
}

type Producers struct {
	TCP  string
	HTTP string
}

func NewClient(conf Config, apiConf *api.Config) Client {
	prod, _ := api.NewProducer(conf.Producers.TCP, apiConf)

	return &client{Config: conf, APIConfig: apiConf, Producer: prod}
}

func (c *client) CreateTopic(name string) *Topic {
	cmd := exec.Command("curl", "-X", "POST", fmt.Sprintf("http://%s/topic/create?topic=%s", c.Producers.HTTP, name))
	cmd.Start()
	cmd.Wait()

	return c.Topic(name)
}

func (c *client) CreateChannel(topic, ch string) *Channel {
	cmd := exec.Command("curl", "-X", "POST", fmt.Sprintf("http://%s/channel/create?topic=%s&channel=%s", c.Producers.HTTP, topic, ch))
	cmd.Start()
	cmd.Wait()

	return c.Channel(topic, ch)
}

func (c *client) NewConsumer(ch *Channel) (*Consumer, error) {
	consumer, err := api.NewConsumer(ch.Topic, ch.Name, c.APIConfig)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Topic:    ch.Topic,
		Channel:  ch.Name,
		Lookupds: c.Lookupds,
		consumer: consumer,
	}, nil
}

func (c *client) Channel(topic, name string) *Channel {
	return &Channel{Topic: topic, Name: name, Producer: c.Producer}
}

func (c *client) Topic(name string) *Topic {
	return &Topic{Name: name, Producer: c.Producer}
}

type Consumer struct {
	Topic    string
	Channel  string
	Lookupds []string
	Lookupd  string
	consumer *api.Consumer
}

func (c *Consumer) AddHandler(handler func(msg *api.Message) error) {
	c.consumer.AddHandler(api.HandlerFunc(func(m *api.Message) error {
		return handler(m)
	}))
}

func (c *Consumer) ConnectToNSQLookupds() {
	c.consumer.ConnectToNSQLookupds(c.Lookupds)
}

type Channel struct {
	Topic    string
	Name     string
	Producer *api.Producer
}

type Topic struct {
	Name     string
	Producer *api.Producer
}

func (t *Topic) Publish(body []byte) error {
	fmt.Println(t.Name)
	return t.Producer.Publish(t.Name, body)
}

func (t *Topic) MultiPublish(body [][]byte) error {
	return t.Producer.MultiPublish(t.Name, body)
}
