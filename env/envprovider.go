package envprovider

import "os"
import "errors"
import "github.com/sohlich/getconfig"
import "fmt"

type envProvider struct {
}

func (c *envProvider) Get(s string) (string, error) {
	out := os.Getenv(s)
	fmt.Printf("envconfig: got value %s for key %s\n", out, s)
	if len(out) == 0 {
		return "", errors.New("envconfig: variable does not exist")
	}
	return out, nil
}

func init() {
	fmt.Println("Setting ENV provider")
	getconfig.RegisterProvider(&envProvider{})
}
