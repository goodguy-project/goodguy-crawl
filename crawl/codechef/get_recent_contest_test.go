package codechef

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/v2/util/jsonx"
)

func TestGetRecentContest(t *testing.T) {
	resp, err := GetRecentContest(nil)
	if err != nil {
		t.Errorf("codechef GetRecentContest failed, err: %v", err)
		return
	}
	t.Logf("codechef GetRecentContest result: %s", jsonx.MarshalString(resp))
}
