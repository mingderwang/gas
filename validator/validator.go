package validator

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	errorMessage map[string]string = map[string]string{
		"NotEmpty": "%s Can not be empty",
	}
)

type Validator struct {
	errorMsg map[string]string
}

func New() *Validator {
	v := &Validator{}
	v.errorMsg = make(map[string]string, 0)

	return v
}

func (v *Validator) Validate(data interface{}, role map[string]string) (bool, error) {
	// d := data.(map[string]interface{})

	// for key, value := range d {
	//     switch role[key] {
	//         case "NotEmpty":
	//             res := NotEmpty(value)
	//             if !res {
	//                 return errors.New(fmt.Sprintf(errorMessage[key], key))
	//             }

	//     }
	// }
	var err error

	rt := reflect.TypeOf(data)

	switch rt.Kind() {
	case reflect.Map:
		err = v.validaeMap(data, role)
	case reflect.Slice:
	default:
	}

	return !v.HasError(), err
}

func (v *Validator) validaeMap(data interface{}, role map[string]string) error {
	switch dt := data.(type) {
	case map[string]string:
		v.validateMapStringString(data.(map[string]string), role)
	case map[string]interface{}:
		v.validateMapStringInterface(data.(map[string]interface{}), role)
	default:
		return errors.New(fmt.Sprint(dt) + " not support to validate")
	}

	return nil
}

func (v *Validator) validateMapStringInterface(data map[string]interface{}, role map[string]string) bool {

	return false
}

func (v *Validator) validateMapStringString(data map[string]string, role map[string]string) {
	for key, value := range role {
		var res bool
		d, ok := data[key]
		if ok {
			res = v.doValidate(d, value)
		} else {
			res = false
		}
		if !res {
			// set error message
			errmsg := fmt.Sprintf(errorMessage[role[key]], key)
			if v.errorMsg == nil {
				v.errorMsg = make(map[string]string, 0)
			}
			v.errorMsg[key] = errmsg
		}
	}

}

func (v *Validator) doValidate(value interface{}, role string) (res bool) {
	switch role {
	case "NotEmpty":
		res = NotEmpty(value)
	}

	return
}

func (v *Validator) HasError() bool {
	return len(v.errorMsg) != 0
}

func (v *Validator) GetErrors() error {
	s := v.GetErrorMessageString()

	return errors.New(s)
}

func (v *Validator) GetErrorMessages() map[string]string {
	return v.errorMsg
}

func (v *Validator) GetErrorMessageByKey(k string) string {
	m, ok := v.errorMsg[k]
	if ok {
		return m
	}

	return k + " is valid(No error)"
}

func (v *Validator) GetErrorMessageString() string {
	var finalmsg string
	for _, msg := range v.errorMsg {
		finalmsg += msg + "\n"
	}

	return finalmsg
}

// Validate functions
func NotEmpty(data interface{}) (res bool) {
	switch data.(type) {
	case nil:
		res = false
	case string:
		s := data.(string)
		if s == "" {
			res = false
		} else {
			res = true
		}

	default:
		res = false
	}

	return
}
