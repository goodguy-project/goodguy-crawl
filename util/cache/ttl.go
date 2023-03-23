package cache

import (
	"container/list"
	"reflect"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

type ttlValue struct {
	Hash   uint64
	Time   time.Time
	Input  []reflect.Value
	Output []reflect.Value
}

type TTLConfig struct {
	TTL     time.Duration
	MaxSize int
}

func TTLWrap[Function any](function Function, config TTLConfig) Function {
	if config.TTL <= 0 {
		config.TTL = time.Hour
	}
	if config.MaxSize <= 0 {
		config.MaxSize = 128
	}
	valueOf := reflect.ValueOf(function)
	typeOf := valueOf.Type()
	if valueOf.Kind() != reflect.Func {
		panic(&reflect.ValueError{Method: "TTLWrap", Kind: valueOf.Kind()})
	}
	hashMap := make(map[uint64]*ttlValue)
	ttlQueue := list.New()
	lock := sync.RWMutex{}
	return reflect.MakeFunc(typeOf, func(input []reflect.Value) []reflect.Value {
		h := hash(input...)
		find := func() (bool, []reflect.Value) {
			if v, ok := hashMap[h]; ok {
				if time.Now().Before(v.Time.Add(config.TTL)) {
					ok := true
					for i, cachedInput := range v.Input {
						x, y := input[i].Interface(), cachedInput.Interface()
						// check deep equal for proto message
						if x, o := x.(proto.Message); o {
							if y, o := y.(proto.Message); o {
								if proto.Equal(x, y) {
									continue
								} else {
									ok = false
									break
								}
							}
						}
						// check deep equal for others
						if !reflect.DeepEqual(x, y) {
							ok = false
							break
						}
					}
					if ok {
						return true, v.Output
					}
				}
			}
			return false, nil
		}
		if ok, op := func() (bool, []reflect.Value) {
			lock.RLock()
			defer lock.RUnlock()
			return find()
		}(); ok {
			return op
		}
		callTime := time.Now()
		output := valueOf.Call(input)
		lock.Lock()
		defer lock.Unlock()
		if ok, op := find(); ok {
			return op
		}
		v := &ttlValue{
			Hash:   h,
			Time:   callTime,
			Input:  input,
			Output: output,
		}
		hashMap[h] = v
		ttlQueue.PushBack(v)
		now := time.Now()
		del := func(v1 *list.Element) {
			v2 := v1.Value.(*ttlValue)
			h := v2.Hash
			if v3, ok := hashMap[h]; ok {
				if v2 == v3 {
					delete(hashMap, h)
				}
			}
			ttlQueue.Remove(v1)
		}
		for true {
			done := false
			if ttlQueue.Len() > 0 {
				front := ttlQueue.Front()
				if now.After(front.Value.(*ttlValue).Time.Add(config.TTL)) {
					del(front)
					done = true
				}
			}
			if !done {
				break
			}
		}
		for ttlQueue.Len() > config.MaxSize {
			del(ttlQueue.Front())
		}
		return output
	}).Interface().(Function)
}
