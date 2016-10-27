package consul

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/hashicorp/consul/api"
	. "gopkg.in/check.v1"
)

type TestConfig struct {
	Host string `consul:"consul_host"`
	Port int    `consul:"consul_port"`
}

func Test(t *testing.T) { TestingT(t) }

type DefTestSuite struct{}

var _ = Suite(&DefTestSuite{})

func (s *DefTestSuite) TestProcess(c *C) {
	cfg := &TestConfig{}

	consul, _ := api.NewClient(api.DefaultConfig())
	Process(cfg, consul.KV())
	c.Assert(cfg.Host, Equals, "10.0.0.1")

	log.Println("### TestProcess result ###")
	out, _ := json.Marshal(cfg)
	log.Printf("%s", out)

}
