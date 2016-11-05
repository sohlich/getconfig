package config

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type ConfigProvider interface {
	Get(s string) (string, error)
}

var (
	provider ConfigProvider
)

func RegisterProvider(p ConfigProvider) {
	if provider != nil {
		log.Panic("Provider already registred")
	}
	provider = p
}

func Process(c interface{}) error {
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
		setField(fVal, fSpec, provider)
	}
	return nil
}

func setField(fVal reflect.Value, fSpec reflect.StructField, prov ConfigProvider) error {
	fName := fSpec.Tag.Get("consul")
	if len(fName) == 0 {
		fName = fSpec.Name
	}
	inVal, err := prov.Get(fName)
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
	return writeField(fVal, fType, inVal)
}

func writeField(fVal reflect.Value, fType reflect.Type, inVal string) error {

	switch fVal.Kind() {
	case reflect.String:
		fVal.SetString(inVal)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)

		if fVal.Kind() == reflect.Int64 && fType.PkgPath() == "time" && fType.Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(inVal)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(inVal, 0, fType.Bits())
		}
		if err != nil {
			return err
		}
		fVal.SetInt(val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(inVal, 0, fType.Bits())
		if err != nil {
			return err
		}
		fVal.SetUint(val)

	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(inVal, fType.Bits())
		if err != nil {
			return err
		}
		fVal.SetFloat(val)

	case reflect.Bool:
		b, err := strconv.ParseBool(inVal)
		if err != nil {
			return err
		}
		fVal.SetBool(b)

	case reflect.Slice:
		vals := strings.Split(inVal, ",")
		sl := reflect.MakeSlice(fType, len(vals), len(vals))
		for i, val := range vals {
			err := writeField(sl.Index(i), fType, val)
			if err != nil {
				return err
			}
		}
		fVal.Set(sl)
	}

	return nil
}
