# getconfig

Parse your config from various providers

```
import "github.com/sohlich/getconfig"
```

As a default provider the env provider could be used.

```
import (
    "github.com/sohlich/getconfig"
    _ "github.com/sohlich/getconfig/env"
    
)

func main(){
        os.Setenv("host", "localhost")
	os.Setenv("port", "22")


	cfg := &struct {
		Port string
		Host string
	}{}

        // Fill the 
	// config struct
	getconfig.Process(cfg)
}

```


Custom providers could be registered like.

```
// Simple map provider
type MapProvider map[string]string

func (m MapProvider) Get(key string) (string, error) {
	iMap := map[string]string(m)
	return iMap[key], nil
}

func main(){


	provider := MapProvider{
		"host":"localhost",
	}

	getconfig.Register(provider)

}

```