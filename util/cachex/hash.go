package cachex

import (
	"math"
	"os"
	"reflect"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func splitmix64(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}

var (
	defaultHash        = splitmix64(splitmix64(uint64(os.Getpid())) ^ splitmix64(uint64(time.Now().Nanosecond())))
	defaultPointerHash = splitmix64(defaultHash)
	defaultStructHash  = splitmix64(defaultPointerHash)
	defaultSliceHash   = splitmix64(defaultStructHash)
	defaultMapHash     = splitmix64(defaultSliceHash)
	defaultBoolHash    = splitmix64(defaultMapHash)
	defaultFloatHash   = splitmix64(defaultBoolHash)
	defaultComplexHash = splitmix64(defaultFloatHash)
)

func hashValueRecursion(isList bool, isMap bool, kind protoreflect.Kind, value protoreflect.Value) uint64 {
	if !isList && kind != protoreflect.MessageKind {
		return hashRecursion(reflect.ValueOf(value.Interface()))
	}
	if isList {
		r := defaultSliceHash
		list := value.List()
		for i := 0; i < list.Len(); i++ {
			r = splitmix64(r + hashValueRecursion(false, false, kind, list.Get(i)))
		}
		return r
	}
	if isMap {
		r := defaultMapHash
		mp := value.Map()
		mp.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			r = r ^ (hashValueRecursion(false, false, kind, protoreflect.Value(k)) + splitmix64(hashValueRecursion(false, false, kind, v)))
			return true
		})
		return r
	}
	return hashMessage(value.Message())
}

func hashMessage(m protoreflect.Message) uint64 {
	r := defaultStructHash
	fields := m.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		value := m.Get(field)
		r = splitmix64(r + hashValueRecursion(field.IsList(), field.IsMap(), field.Kind(), value))
	}
	return r
}

func hashRecursion(x reflect.Value) uint64 {
	if x.CanInterface() {
		if m, ok := x.Interface().(protoreflect.ProtoMessage); ok {
			return hashMessage(m.ProtoReflect())
		}
	}
	switch x.Kind() {
	case reflect.Pointer:
		if !x.IsValid() {
			return defaultPointerHash
		}
		return hashRecursion(x.Elem())
	case reflect.Struct:
		if x.CanInterface() {
			if t, ok := x.Interface().(time.Time); ok {
				return splitmix64(uint64(t.UnixNano()))
			}
		}
		r := defaultStructHash
		for i := 0; i < x.NumField(); i++ {
			r = splitmix64(r + hashRecursion(x.Field(i)))
		}
		return r
	case reflect.Array, reflect.Slice:
		if x.Kind() == reflect.Slice && x.IsNil() {
			return defaultSliceHash
		}
		r := defaultSliceHash
		for i := 0; i < x.Len(); i++ {
			r = splitmix64(r + hashRecursion(x.Index(i)))
		}
		return r
	case reflect.Map:
		if x.IsNil() {
			return defaultMapHash
		}
		r := defaultMapHash
		for _, k := range x.MapKeys() {
			v := x.MapIndex(k)
			r = r ^ (hashRecursion(k) + splitmix64(hashRecursion(v)))
		}
		return r
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return splitmix64(defaultHash + uint64(x.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return splitmix64(defaultHash + x.Uint())
	case reflect.Bool:
		if x.Bool() {
			return splitmix64(defaultBoolHash + 1)
		}
		return splitmix64(defaultBoolHash)
	case reflect.Float32, reflect.Float64:
		return splitmix64(defaultFloatHash + math.Float64bits(x.Float()))
	case reflect.Complex64, reflect.Complex128:
		y := x.Complex()
		a := math.Float64bits(real(y))
		b := math.Float64bits(imag(y))
		return splitmix64(splitmix64(defaultComplexHash+a) + b)
	case reflect.String:
		return hashRecursion(reflect.ValueOf([]byte(x.String())))
	default:
		panic(&reflect.ValueError{Method: "Hash", Kind: x.Kind()})
	}
}

func Hash(x ...reflect.Value) uint64 {
	r := defaultSliceHash
	for _, y := range x {
		r = splitmix64(r + hashRecursion(y))
	}
	return r
}
