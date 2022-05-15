package reqbind

import (
	"encoding"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var (
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

func Bind(target interface{}, r *http.Request) error {
	if err := doBind(target, r); err != nil {
		return err
	}

	if err := doValidate(target, r); err != nil {
		return err
	}

	return nil
}

func doBind(target interface{}, r *http.Request) error {
	if r.Header.Get("Content-type") == "application/json" {
		// JSON body
		if err := json.NewDecoder(r.Body).Decode(target); err != nil {
			return err
		}
	}

	return doFormBind(target, r)
}

func doFormBind(target interface{}, r *http.Request) error {
	v := reflect.ValueOf(target)
	if (v.Kind() != reflect.Ptr) || (v.Elem().Kind() != reflect.Struct) {
		return errors.New("target must be a pointer to a struct")
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	return bindStruct(v.Elem(), r.Form, "")
}

func bindStruct(sct reflect.Value, values url.Values, prefix string) error {
	sctType := sct.Type()
	for i := 0; i < sctType.NumField(); i++ {
		fieldName := sctType.Field(i)

		urlTag, ok := fieldName.Tag.Lookup("req")
		if !ok {
			continue
		}

		field := sct.FieldByName(fieldName.Name)

		formName, option, hasOption := strings.Cut(urlTag, ",")
		if hasOption && option == "zero" {
			field.Set(reflect.Zero(field.Type()))
		}
		fullFormName := prefix + formName

		var err error
		switch field.Type().Kind() {
		case reflect.Struct:
			err = bindStruct(field, values, fullFormName+".")
		default:
			if !values.Has(fullFormName) {
				continue
			}

			// Use the last form value for scalar types.  This is to deal with the stupid checkbox hack
			lastValue := values[fullFormName][len(values[fullFormName])-1]
			err = errors.Wrapf(setScalar(field, lastValue), "field %v", fullFormName)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func setScalar(field reflect.Value, formValue string) error {
	// Primitives
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(formValue)
		return nil
	case reflect.Int:
		intValue, _ := strconv.Atoi(formValue)
		field.SetInt(int64(intValue))
		return nil
	case reflect.Bool:
		switch formValue {
		case "1", "t", "T", "true", "TRUE", "True", "on", "ON":
			field.SetBool(true)
		case "0", "f", "F", "false", "FALSE", "False", "off", "OFF":
			field.SetBool(false)
		}
		return nil
	}

	// Interfaces that support text unmarshalling
	if field.Type().AssignableTo(textUnmarshalerType) {
		ut := field.Interface().(encoding.TextUnmarshaler)
		_ = ut.UnmarshalText([]byte(formValue))
		return nil
	} else if fieldPtr := field.Addr(); fieldPtr.Type().AssignableTo(textUnmarshalerType) {
		ut := fieldPtr.Interface().(encoding.TextUnmarshaler)
		_ = ut.UnmarshalText([]byte(formValue))
		return nil
	}

	return errors.New("unsupported type")
}
