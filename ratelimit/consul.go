package ratelimit

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/ulule/limiter"
)

type ConsulStore struct {
	client *api.Client
	kv     *api.KV
	prefix string
}

func NewConsulStore(c *api.Client) limiter.Store {
	return &ConsulStore{
		client: c,
		kv:     c.KV(),
		prefix: "RateLimit",
	}
}

func (s *ConsulStore) Get(key string, rate limiter.Rate) (limiter.Context, error) {
	//ctx := limiter.Context{}
	key = fmt.Sprintf("%s/%s", s.prefix, key)
	ms := int64(time.Millisecond)
	now := time.Now()
	exp := (now.UnixNano()/ms + int64(rate.Period)/ms) / 1000
	pair, _, _ := s.kv.Get(key, &api.QueryOptions{})
	count, expiration, increased := s.getValue(pair)

	if pair == nil || s.expired(expiration) {
		s.kv.Put(&api.KVPair{
			Key:   key,
			Value: []byte(fmt.Sprintf("%d:%d", 1, exp)),
		}, &api.WriteOptions{})

		return limiter.Context{
			Limit:     rate.Limit,
			Remaining: rate.Limit - 1,
			Reset:     exp,
			Reached:   false,
		}, nil
	}

	// update
	s.kv.Put(&api.KVPair{
		Key:   key,
		Value: []byte(fmt.Sprintf("%s:%d", increased, exp)),
	}, &api.WriteOptions{})

	return s.getContextFromState(now, rate, expiration, count), nil
}

func (s *ConsulStore) Peek(key string, rate limiter.Rate) (limiter.Context, error) {
	//ctx := limiter.Context{}
	key = fmt.Sprintf("%s/%s", s.prefix, key)
	ms := int64(time.Millisecond)
	now := time.Now()
	exp := (now.UnixNano()/ms + int64(rate.Period)/ms) / 1000
	pair, _, _ := s.kv.Get(key, &api.QueryOptions{})
	count, expiration, _ := s.getValue(pair)

	if pair == nil || s.expired(expiration) {
		return limiter.Context{
			Limit:     rate.Limit,
			Remaining: rate.Limit,
			Reset:     exp,
			Reached:   false,
		}, nil
	}

	return s.getContextFromState(now, rate, expiration, count), nil
}

//
func (s *ConsulStore) getValue(pair *api.KVPair) (int64, int64, []byte) {
	if pair == nil {
		return 0, 0, []byte("0")
	}

	val := pair.Value
	str := strings.Split(string(val), ":")
	count, expiration := str[0], str[1]

	c, _ := strconv.Atoi(count)
	e, _ := strconv.Atoi(expiration)

	return int64(c), int64(e), []byte(strconv.Itoa(c + 1))
}

//
func (s *ConsulStore) expired(t int64) bool {
	if time.Now().Unix() > t {
		return true
	}

	return false
}

func (s *ConsulStore) getContextFromState(now time.Time, rate limiter.Rate, expiration, count int64) limiter.Context {
	remaining := int64(0)
	if count < rate.Limit {
		remaining = rate.Limit - count
	}

	expire := time.Unix(0, expiration)

	return limiter.Context{
		Limit:     rate.Limit,
		Remaining: remaining,
		Reset:     expire.Add(time.Duration(expire.Sub(now).Seconds()) * time.Second).Unix(),
		Reached:   count > rate.Limit,
	}
}
