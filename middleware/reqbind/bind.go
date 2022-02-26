package reqbind

import (
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"strconv"
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
	v := reflect.ValueOf(target)
	if (v.Kind() != reflect.Ptr) || (v.Elem().Kind() != reflect.Struct) {
		return errors.New("target must be a pointer to a struct")
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	sct := v.Elem()
	sctType := sct.Type()
	for i := 0; i < sctType.NumField(); i++ {
		fieldName := sctType.Field(i)
		urlTag, ok := fieldName.Tag.Lookup("req")
		if !ok {
			continue
		}

		field := sct.FieldByName(fieldName.Name)
		value := r.FormValue(urlTag)
		if err := setField(field, value); err != nil {
			return nil
		}
	}

	return nil
}

func setField(field reflect.Value, formValue string) error {
	switch field.Type().Kind() {
	case reflect.String:
		field.Set(reflect.ValueOf(formValue))
	case reflect.Int:
		intValue, _ := strconv.Atoi(formValue)
		field.Set(reflect.ValueOf(intValue))
	}

	return nil
}
