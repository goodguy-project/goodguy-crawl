package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/goodguy-project/goodguy-crawl/proto"
)

func TestTTLWrap(t *testing.T) {
	f1 := func(x string) string {
		time.Sleep(2 * time.Second)
		return "ok"
	}
	f2 := TTLWrap(f1, TTLConfig{
		MaxSize: 3,
	})
	start := time.Now()
	f2("2")
	f2("3")
	f2("2")
	end := time.Now()
	fmt.Println(end.Sub(start))
}

func TestProtoTTLWrap(t *testing.T) {
	f1 := func(*proto.GetRecentContestRequest) string {
		time.Sleep(1 * time.Second)
		return "ok"
	}
	f2 := TTLWrap(f1, TTLConfig{
		MaxSize: 3,
	})
	start := time.Now()
	a := &proto.GetRecentContestRequest{}
	a.Reset()
	a.Platform = "2"
	b := &proto.GetRecentContestRequest{
		Platform: "3",
	}
	c := &proto.GetRecentContestRequest{
		Platform: "2",
	}
	f2(a)
	f2(b)
	f2(c)
	end := time.Now()
	fmt.Println(end.Sub(start))
}
