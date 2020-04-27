package eren

import (
	"encoding/json"
	"reflect"
)
const (
	ErrorTypeBind    ErrorType = 1 << 63 // used when c.Bind() fails
	ErrorTypeRender  ErrorType = 1 << 62 // used when c.Render() fails
	ErrorTypePrivate ErrorType = 1 << 0  // used by custom
	ErrorTypePublic  ErrorType = 1 << 1

	ErrorTypeAny ErrorType = 1<<64 - 1
	ErrorTypeNu            = 2
)

type ErrorType uint64
type errors []*Error

type Error struct {
	Err error
	Type ErrorType
	Meta interface{}
}

func (e Error) Error() string {
	return e.Err.Error()
}

var _ error = &Error{}

func (e *Error) SetType(etype ErrorType) *Error{
	e.Type = etype
	return e
}

func (e *Error) SetMeta(data interface{}) *Error {
	e.Meta = data
	return e
}

func (e *Error) JSON() interface{} {
	j := map[string]interface{}{}
	if e.Meta != nil {
		value := reflect.ValueOf(e.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return e.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				j[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			j["meta"] = e.Meta
		}
	}
	if _, ok := j["error"]; !ok {
		j["error"] = e.Error()
	}
	return j
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.JSON())
}

