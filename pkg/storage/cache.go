package storage

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type Template struct {
	LastName string `json:"last_name"`
	Birthday string `json:"birthday"`
}

type Cache struct {
	client *memcache.Client
}

func NewMemcached() (*Cache, error) {
	// XXX Assuming environment variable contains only one server
	client := memcache.New("localhost:11211")

	if err := client.Ping(); err != nil {
		return nil, err
	}

	client.Timeout = 100 * time.Millisecond
	client.MaxIdleConns = 100

	return &Cache{
		client: client,
	}, nil
}

func (c *Cache) FetchTemplate(reqId string) (*Template, error) {
	item, err := c.client.Get(reqId)
	if err != nil {
		return &Template{}, err
	}

	b := bytes.NewReader(item.Value)

	var res Template

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return &Template{}, err
	}

	return &res, nil
}

func (c *Cache) CacheTemplate(t *Template) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(*t); err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:        t.LastName,
		Value:      b.Bytes(),
		Expiration: int32(time.Now().Add(25 * time.Second).Unix()),
	})
}
