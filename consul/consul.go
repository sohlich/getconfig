package consul

import (
	con "github.com/hashicorp/consul/api"
	"github.com/sohlich/consulconfig"
)

type consulProvider struct {
	client *con.KV
}

func (c *consulProvider) Get(s string) (string, error) {
	p, _, err := c.client.Get(s, nil)
	if p != nil {
		return string(p.Value), nil
	}
	return "", err
}

func Process(c interface{}, kv *con.KV) error {
	return config.Process(c, &consulProvider{kv})
}
