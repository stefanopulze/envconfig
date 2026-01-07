package envconfig

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/stefanopulze/envconfig/internal/dotenv"
)

const (
	defaultSeparator = ","
	tagEnvSeparator  = "env-separator"
	tagEnvDefault    = "env-default"
	tagEnv           = "env"
	tagEnvPrefix     = "env-prefix"
)

type Setter interface {
	SetValue(string) error
}

func ReadEnv(cfg interface{}) error {
	s := reflect.ValueOf(cfg)

	// unwrap pointer
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	// process-only structures
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("wrong type %v", s.Kind())
	}

	return readStruct(s, "")
}

func ReadDotEnv(filename string) error {
	return dotenv.Parse(filename)
}

func readStruct(s reflect.Value, prefix string) error {
	// unwrap pointer
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	var err error
	for i := 0; i < s.NumField(); i++ {
		fieldType := s.Type().Field(i)
		fieldValue := s.Field(i)

		if fieldType.Type.Kind() == reflect.Struct {
			cf := joinPrefix(prefix, lookupEnvPrefix(fieldType))
			if err = readStruct(fieldValue, cf); err != nil {
				return err
			}
			continue
		}

		envName := joinPrefix(prefix, lookupEnvName(fieldType))
		value, found := os.LookupEnv(envName)
		if !found {
			defaultValue := fieldType.Tag.Get(tagEnvDefault)
			if len(defaultValue) > 0 {
				value = defaultValue
			}
		}

		err = parseValue(fieldType, fieldValue, value)

		if err != nil {
			return errors.Join(errors.New(fmt.Sprintf("cannot convert value or missing env-default for field: %s", envName)), err)
		}
	}

	return nil
}

func parseValue(fieldType reflect.StructField, field reflect.Value, value string) error {
	if field.CanInterface() {
		if ct, ok := field.Interface().(encoding.TextUnmarshaler); ok {
			return ct.UnmarshalText([]byte(value))
		} else if ctp, ok := field.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return ctp.UnmarshalText([]byte(value))
		}

		if cs, ok := field.Interface().(Setter); ok {
			return cs.SetValue(value)
		} else if csp, ok := field.Addr().Interface().(Setter); ok {
			return csp.SetValue(value)
		}
	}

	valueType := field.Type()

	switch valueType.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		number, err := strconv.ParseInt(value, 0, valueType.Bits())
		if err != nil {
			return err
		}
		field.SetInt(number)

	case reflect.Int64:
		if valueType == reflect.TypeOf(time.Duration(0)) {
			// try to parse time
			d, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			field.SetInt(int64(d))
		} else {
			// parse regular integer
			number, err := strconv.ParseInt(value, 0, valueType.Bits())
			if err != nil {
				return err
			}
			field.SetInt(number)
		}

		// parse unsigned integer value
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		number, err := strconv.ParseUint(value, 0, valueType.Bits())
		if err != nil {
			return err
		}
		field.SetUint(number)

	// parse floating point value
	case reflect.Float32, reflect.Float64:
		number, err := strconv.ParseFloat(value, valueType.Bits())
		if err != nil {
			return err
		}
		field.SetFloat(number)

	// parse sliced value
	case reflect.Slice:
		separator := lookupEnvSeparator(fieldType)
		sliceValue, err := parseSlice(fieldType, valueType, value, separator)
		if err != nil {
			return err
		}

		field.Set(*sliceValue)

	// parse mapped value
	case reflect.Map:
		separator := lookupEnvSeparator(fieldType)
		mapValue, err := parseMap(fieldType, valueType, value, separator)
		if err != nil {
			return err
		}

		field.Set(*mapValue)

	default:
		return fmt.Errorf("unsupported type %s.%s", valueType.PkgPath(), valueType.Name())
	}

	return nil
}

func parseMap(fType reflect.StructField, valueType reflect.Type, value string, sep string) (*reflect.Value, error) {
	mapValue := reflect.MakeMap(valueType)
	if len(strings.TrimSpace(value)) != 0 {
		pairs := strings.Split(value, sep)
		for _, pair := range pairs {
			kvPair := strings.SplitN(pair, ":", 2)
			if len(kvPair) != 2 {
				return nil, fmt.Errorf("invalid map item: %q", pair)
			}
			k := reflect.New(valueType.Key()).Elem()
			err := parseValue(fType, k, kvPair[0])
			if err != nil {
				return nil, err
			}
			v := reflect.New(valueType.Elem()).Elem()
			err = parseValue(fType, v, kvPair[1])
			if err != nil {
				return nil, err
			}
			mapValue.SetMapIndex(k, v)
		}
	}
	return &mapValue, nil
}

func parseSlice(fType reflect.StructField, valueType reflect.Type, value string, sep string) (*reflect.Value, error) {
	sliceValue := reflect.MakeSlice(valueType, 0, 0)

	if valueType.Elem().Kind() == reflect.Uint8 {
		sliceValue = reflect.ValueOf([]byte(value))
	} else if len(strings.TrimSpace(value)) != 0 {
		values := strings.Split(value, sep)
		sliceValue = reflect.MakeSlice(valueType, len(values), len(values))
		for i, val := range values {
			if err := parseValue(fType, sliceValue.Index(i), val); err != nil {
				return nil, err
			}
		}
	}
	return &sliceValue, nil
}

func lookupEnvSeparator(fieldType reflect.StructField) string {
	sep := fieldType.Tag.Get(tagEnvSeparator)
	if len(sep) == 0 {
		return defaultSeparator
	}

	return sep
}

func lookupEnvName(fieldType reflect.StructField) string {
	name := fieldType.Tag.Get(tagEnv)
	if len(name) == 0 {
		name = fieldType.Name
	}

	return strings.ToUpper(name)
}

func lookupEnvPrefix(fieldType reflect.StructField) string {
	name := fieldType.Tag.Get(tagEnvPrefix)
	if len(name) == 0 {
		return ""
	}

	return name
}

func joinPrefix(parent string, current string) string {
	if len(parent) == 0 {
		return current
	} else if len(current) == 0 {
		return parent
	}

	return fmt.Sprintf("%s_%s", parent, current)
}
