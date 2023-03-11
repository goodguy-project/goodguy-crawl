package leetcode

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func TestGetDailyQuestion(t *testing.T) {
	resp, err := GetDailyQuestion(nil)
	if err != nil {
		t.Errorf("leetcode GetDailyQuestion failed, err: %v", err)
		return
	}
	t.Logf("leetcode GetDailyQuestion result: %s", jsonx.MarshalString(resp))
}
