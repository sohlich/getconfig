package config

import (
	"encoding/json"
	"log"
	"testing"

	. "gopkg.in/check.v1"
)

type TestConfig struct {
	Host string `consul:"consul_host"`
	Port int    `consul:"consul_port"`
}

func Test(t *testing.T) { TestingT(t) }

type DefTestSuite struct{}

var _ = Suite(&DefTestSuite{})

type MapProvider map[string]string

func (m MapProvider) Get(key string) (string, error) {
	iMap := map[string]string(m)
	return iMap[key], nil
}

func (s *DefTestSuite) TestProcess(c *C) {
	cfg := &TestConfig{}

	p := MapProvider{
		"consul_host": "10.0.0.1",
		"consul_port": "8080",
	}

	Register(p)
	Process(cfg)
	c.Assert(cfg.Host, Equals, "10.0.0.1")

	log.Println("### TestProcess result ###")
	out, _ := json.Marshal(cfg)
	log.Printf("%s", out)

}
