package envprovider

import (
	"os"
	"testing"

	"fmt"

	"github.com/sohlich/getconfig"
)

func TestEnv(t *testing.T) {
	os.Setenv("host", "localhost")
	os.Setenv("port", "22")
	cfg := &struct {
		Port string
		Host string
	}{}
	getconfig.Process(cfg)

	fmt.Println(cfg)

	if cfg.Host != "localhost" {
		t.Fail()
	}

}
