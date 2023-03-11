package jsonx

import (
	"reflect"

	"github.com/bytedance/sonic"
)

func Marshal(val any) []byte {
	r, _ := sonic.Marshal(val)
	return r
}

func MarshalString(val any) string {
	r, _ := sonic.MarshalString(val)
	return r
}

func Unmarshal[Resp any, Req string | []byte](s Req) (resp Resp, err error) {
	if typeof := reflect.TypeOf(resp); typeof.Kind() == reflect.Pointer {
		resp = reflect.New(typeof.Elem()).Interface().(Resp)
		switch s := any(s).(type) {
		case string:
			err = sonic.UnmarshalString(s, resp)
		case []byte:
			err = sonic.Unmarshal(s, resp)
		}
		return
	}
	switch s := any(s).(type) {
	case string:
		err = sonic.UnmarshalString(s, &resp)
	case []byte:
		err = sonic.Unmarshal(s, &resp)
	}
	return
}
