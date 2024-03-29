package cachex

import (
	"reflect"
	"testing"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
)

func TestHash(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		if Hash(reflect.ValueOf("123")) != Hash(reflect.ValueOf("123")) {
			t.Error("string hash failed")
			return
		}
	})
	t.Run("array", func(t *testing.T) {
		a := [4]int{123, 4235, 65, 677}
		b := a[:]
		if Hash(reflect.ValueOf(a)) != Hash(reflect.ValueOf(b)) {
			t.Error("array hash failed")
			return
		}
	})
	t.Run("map", func(t *testing.T) {
		a := map[string]string{"123": "345", "234": "345"}
		b := map[string]string{"234": "345", "123": "345"}
		if Hash(reflect.ValueOf(a)) != Hash(reflect.ValueOf(b)) {
			t.Error("map hash failed")
			return
		}
	})
	t.Run("struct", func(t *testing.T) {
		a := struct{ A int }{A: 12}
		b := struct{ A int64 }{A: 12}
		if Hash(reflect.ValueOf(a)) != Hash(reflect.ValueOf(b)) {
			t.Error("struct hash failed")
			return
		}
	})
	t.Run("proto", func(t *testing.T) {
		a := &proto.GetRecentContestResponse{}
		a.Reset()
		a.Platform = "1"
		a.RecentContest = []*proto.GetRecentContestResponse_Contest{
			{Name: "1"},
			{Name: "2"},
		}
		b := &proto.GetRecentContestResponse{}
		b.Reset()
		b.Platform = "1"
		b.RecentContest = []*proto.GetRecentContestResponse_Contest{
			{Name: "1"},
			{Name: "2"},
		}
		c := &proto.GetRecentContestResponse{
			Platform: "1",
			RecentContest: []*proto.GetRecentContestResponse_Contest{
				{Name: "2"},
				{Name: "2"},
			}}
		if Hash(reflect.ValueOf(a)) != Hash(reflect.ValueOf(b)) {
			t.Error("proto hash failed")
			return
		}
		if Hash(reflect.ValueOf(a)) == Hash(reflect.ValueOf(c)) {
			t.Error("proto hash failed")
			return
		}
	})
}
