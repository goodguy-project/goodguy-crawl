package jsonx

import (
	"reflect"

	"github.com/bytedance/sonic"
)

var (
	encoder = sonic.Config{}.Froze()
)

func Marshal(val any) []byte {
	r, _ := encoder.Marshal(val)
	return r
}

func MarshalString(val any) string {
	r, _ := encoder.MarshalToString(val)
	return r
}

var (
	decoder = sonic.Config{UseInt64: true}.Froze()
)

func Unmarshal[Resp any, Req string | []byte](s Req) (resp Resp, err error) {
	if typeof := reflect.TypeOf(resp); typeof.Kind() == reflect.Pointer {
		resp = reflect.New(typeof.Elem()).Interface().(Resp)
		switch s := any(s).(type) {
		case string:
			err = decoder.UnmarshalFromString(s, resp)
		case []byte:
			err = decoder.Unmarshal(s, resp)
		}
		return
	}
	switch s := any(s).(type) {
	case string:
		err = decoder.UnmarshalFromString(s, &resp)
	case []byte:
		err = decoder.Unmarshal(s, &resp)
	}
	return
}
