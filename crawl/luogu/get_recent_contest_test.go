package luogu

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func TestGetRecentContest(t *testing.T) {
	resp, err := GetRecentContest(nil)
	if err != nil {
		t.Errorf("luogu GetRecentContest failed, err: %v", err)
		return
	}
	t.Logf("luogu GetRecentContest result: %s", jsonx.MarshalString(resp))
}
