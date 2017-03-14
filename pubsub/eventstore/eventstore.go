package eventstore

import (
	"github.com/jetbasrawi/go.geteventstore"
)

type Config struct {
	URL      string
	Username string
	Password string
}

type Client struct {
	Config Config
	Client *goes.Client
}

func NewClient(conf Config) *Client {
	client, _ := goes.NewClient(nil, conf.URL)
	client.SetBasicAuth(conf.Username, conf.Password)

	return &Client{
		Config: conf,
		Client: client,
	}
}

func (c *Client) CreateStream(name string) *Stream {
	return &Stream{
		Name:   name,
		Reader: c.Client.NewStreamReader(name),
		Writer: c.Client.NewStreamWriter(name),
	}
}

func (c *Client) Stream(name string) *Stream {
	return &Stream{Name: name}
}

type Stream struct {
	Name   string
	Reader *goes.StreamReader
	Writer *goes.StreamWriter
}

type Event struct {
	ID   string
	Type string
	Data interface{}
	Meta interface{}
}

func (s *Stream) AddHandler(handler func(*goes.StreamReader)) {
	go func() {
		for s.Reader.Next() {
			if s.Reader.Err() != nil {
				panic(s.Reader.Err())
			} else {
				handler(s.Reader)
			}
		}
	}()
}

func (s *Stream) Append(e Event) error {
	return s.Writer.Append(nil, goes.NewEvent(e.ID, e.Type, e.Data, e.Meta))
}

func (s *Stream) MultiAppend(e []Event) error {
	var events []*goes.Event
	for _, v := range e {
		events = append(events, goes.NewEvent(v.ID, v.Type, v.Data, v.Meta))
	}
	return s.Writer.Append(nil, events...)
}
