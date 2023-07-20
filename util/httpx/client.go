package httpx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"

	"github.com/bytedance/sonic"
	"golang.org/x/sync/semaphore"
)

var (
	controller = make(map[string]*semaphore.Weighted)
	lock       sync.RWMutex
)

func getSemaphore(key string) *semaphore.Weighted {
	if s := func() *semaphore.Weighted {
		lock.RLock()
		defer lock.RUnlock()
		s, ok := controller[key]
		if ok {
			return s
		}
		return nil
	}(); s != nil {
		return s
	}
	lock.Lock()
	defer lock.Unlock()
	s, ok := controller[key]
	if ok {
		return s
	}
	s = semaphore.NewWeighted(1)
	controller[key] = s
	return s
}

type Response struct {
	Header  http.Header
	Cookies []*http.Cookie
}

func SendRequest[Resp any](key string, client *http.Client, req *http.Request) (resp Resp, extra *Response, err error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	// prepare
	extra = new(Response)
	isPointer := false
	if typeof := reflect.TypeOf(resp); typeof.Kind() == reflect.Pointer {
		resp = reflect.New(typeof.Elem()).Interface().(Resp)
		isPointer = true
	} else if typeof.Kind() == reflect.Slice {
		reflect.ValueOf(&resp).Elem().Set(reflect.MakeSlice(typeof, 0, 0))
	}
	// acquire semaphore
	s := getSemaphore(key)
	err = s.Acquire(context.Background(), 1)
	if err != nil {
		return
	}
	defer s.Release(1)
	// do send request
	var httpResponse *http.Response
	if client == nil {
		client = http.DefaultClient
	}
	httpResponse, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		e := httpResponse.Body.Close()
		if err == nil {
			err = e
		}
	}()
	// handle response
	extra.Cookies = httpResponse.Cookies()
	extra.Header = httpResponse.Header
	var body []byte
	body, err = io.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}
	if len(body) > 0 {
		if _, ok := any(resp).([]byte); ok {
			reflect.ValueOf(&resp).Elem().Set(reflect.ValueOf(body))
		} else if _, ok := any(resp).(struct{}); !ok {
			if isPointer {
				err = sonic.Unmarshal(body, resp)
				if err != nil {
					return
				}
			} else {
				err = sonic.Unmarshal(body, &resp)
				if err != nil {
					return
				}
			}
		}
	}
	if httpResponse.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("error with status code: %d", httpResponse.StatusCode))
		return
	}
	return
}
