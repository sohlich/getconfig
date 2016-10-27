package config

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

type ConfigProvider interface {
	Get(s string) (string, error)
}

func Process(c interface{}, prov ConfigProvider) error {
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

func setField(fVal reflect.Value, fSpec reflect.StructField, prov ConfigProvider) error {
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
		i, err := strconv.ParseInt(provVal, 0, fType.Bits())
		if err != nil {
			return err
		}
		fVal.SetInt(i)

	case reflect.String:
		fVal.SetString(provVal)

	case reflect.Bool:
		b, err := strconv.ParseBool(provVal)
		if err != nil {
			return err
		}
		fVal.SetBool(b)
	}

	return nil
}
