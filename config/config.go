package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func GetConfig(c interface{}, cf string) error {
	if reflect.TypeOf(c).Kind() != reflect.Ptr {
		return errors.New("config.GetConfig() expects a pointer arg")
	}

	// read config yaml file
	raw, err := ioutil.ReadFile(cf)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(raw, c)
	if err != nil {
		return err
	}

	// read env vars
	err = updateEnvFields(reflect.ValueOf(c), "")
	if err != nil {
		return err
	}

	return nil
}

func updateEnvFields(v reflect.Value, prefix string) error {
	if v.Kind() != reflect.Ptr {
		return errors.New("Not a pointer value")
	}

	v = reflect.Indirect(v)

	switch v.Kind() {
	case reflect.Int:
		val := os.Getenv(prefix)
		if val != "" {
			conv, err := strconv.Atoi(val)
			if err == nil {
				v.SetInt(int64(conv))
			}
		}
	case reflect.Float64:
		val := os.Getenv(prefix)
		if val != "" {
			conv, err := strconv.ParseFloat(val, 64)
			if err == nil {
				v.SetFloat(conv)
			}
		}
	case reflect.String:
		val := os.Getenv(prefix)
		if val != "" {
			v.SetString(val)
		}
	case reflect.Bool:
		val := os.Getenv(prefix)
		if val != "" {
			conv, err := strconv.ParseBool(val)
			if err == nil {
				v.SetBool(conv)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			vi := v.Index(i)
			name := strconv.Itoa(i)
			if prefix != "" {
				name = prefix + "_" + name
			}
			err := updateEnvFields(vi.Addr(), name)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		vt := reflect.TypeOf(v.Interface())
		for i := 0; i < vt.NumField(); i++ {
			ft := vt.Field(i)
			fv := v.Field(i)
			name := strings.ToUpper(ft.Name)
			if prefix != "" {
				name = prefix + "_" + name
			}
			err := updateEnvFields(fv.Addr(), name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
