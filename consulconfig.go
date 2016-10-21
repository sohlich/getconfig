package consulconfig

import (
	"fmt"
	"reflect"
	"strconv"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type configProvider interface {
	Get(s string) (string, error)
}

type consulProvider struct {
	client *consul.KV
}

func (c *consulProvider) Get(s string) (string, error) {
	p, _, err := c.client.Get(s, nil)
	if p != nil {
		return string(p.Value), nil
	}
	return "", err
}

func Process(c interface{}, kv *consul.KV) error {
	return process(c, &consulProvider{kv})
}

func process(c interface{}, prov configProvider) error {
	t := reflect.ValueOf(c)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("Config cannot be pointer")
	}
	t = t.Elem()
	fieldCount := t.NumField()
	typeOfSpec := t.Type()
	for idx := 0; idx < fieldCount; idx++ {
		fVal := t.Field(idx)
		fSpec := typeOfSpec.Field(idx)
		setField(fVal, fSpec, prov)
	}
	return nil
}

func setField(fVal reflect.Value, fSpec reflect.StructField, prov configProvider) error {
	fName := fSpec.Tag.Get("consul")
	if len(fName) == 0 {
		fName = fSpec.Name
	}
	provVal, err := prov.Get(fName)
	if err != nil {
		errors.Wrap(err, "Cannot parse value from consul")
	}

	fType := fSpec.Type

	if !fVal.IsValid() {
		return fmt.Errorf("Cannot set field")
	}
	if !fVal.CanSet() {
		return fmt.Errorf("Cannot set field")
	}

	switch fVal.Kind() {
	case reflect.Int:
		strconv.ParseInt(provVal, 0, fType.Bits())

	case reflect.String:
		fVal.SetString(provVal)
	}

	return nil
}
