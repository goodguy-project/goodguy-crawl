package acwing

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/v2/util/jsonx"
)

func TestGetRecentContest(t *testing.T) {
	resp, err := GetRecentContest(nil)
	if err != nil {
		t.Errorf("acwing GetRecentContest failed, err: %v", err)
		return
	}
	t.Logf("acwing GetRecentContest result: %s", jsonx.MarshalString(resp))
}
