package leetcode

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func TestGetRecentContest(t *testing.T) {
	resp, err := GetRecentContest(nil)
	if err != nil {
		t.Errorf("leetcode GetRecentContest failed, err: %v", err)
		return
	}
	t.Logf("leetcode GetRecentContest result: %s", jsonx.MarshalString(resp))
}
