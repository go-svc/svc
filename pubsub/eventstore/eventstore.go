package eventstore

import (
	"github.com/jetbasrawi/go.geteventstore"
)

// Config 是 EventStore 建立客戶端時的設定檔。
type Config struct {
	// URL 是連線到 EventStore 的 HTTP 網址。
	URL      string
	Username string
	Password string
}

// Client 是 EventStore 的客戶端。
type Client struct {
	Config Config
	Client *goes.Client
}

// NewClient 會建立一個新的客戶端。
func NewClient(conf Config) *Client {
	client, _ := goes.NewClient(nil, conf.URL)
	client.SetBasicAuth(conf.Username, conf.Password)

	return &Client{Config: conf, Client: client}
}

// CreateStream 會建立並返回一個串流。
func (c *Client) CreateStream(name string) *Stream {
	return c.Stream(name)
}

// Stream 會回傳一個指定的串流。
func (c *Client) Stream(name string) *Stream {
	return &Stream{
		Name:   name,
		Reader: c.Client.NewStreamReader(name),
		Writer: c.Client.NewStreamWriter(name),
	}
}

// Stream 呈現了一個事件串流。
type Stream struct {
	// Name 是這個串流的名稱。
	Name   string
	Reader *goes.StreamReader
	Writer *goes.StreamWriter
}

// Event 呈現了一個事件資料。
type Event struct {
	// ID 是這個事件的編號，若留空則會以 UUID 自動填上。
	ID string
	// Type 是這個事件的自訂種類。
	Type string
	// Data 是這個事件帶有的資料。
	Data interface{}
	// Meta 是事件的中繼資料，這裡可以放置一些額外的資料像是事件的產生者是誰⋯⋯等。
	Meta interface{}
}

// AddHandler 會替這個串流新增一個處理函式，一但接收到任何事件就會呼叫這個函式。
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

// Append 會在此串流上推送一個新的事件。
func (s *Stream) Append(e Event) error {
	return s.Writer.Append(nil, goes.NewEvent(e.ID, e.Type, e.Data, e.Meta))
}

// MultiAppend 會在此串流上一次推送多個事件。
func (s *Stream) MultiAppend(e []Event) error {
	var events []*goes.Event
	for _, v := range e {
		events = append(events, goes.NewEvent(v.ID, v.Type, v.Data, v.Meta))
	}

	return s.Writer.Append(nil, events...)
}
