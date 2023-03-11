package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestTTLWrap(t *testing.T) {
	f1 := func(x string) string {
		time.Sleep(2 * time.Second)
		return "ok"
	}
	f2 := TTLWrap(f1, TTLConfig{
		MaxSize: 1,
	})
	start := time.Now()
	f2("2")
	f2("3")
	f2("2")
	end := time.Now()
	fmt.Println(end.Sub(start))
}
