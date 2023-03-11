package nowcoder

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func TestGetRecentContest(t *testing.T) {
	resp, err := GetRecentContest(nil)
	if err != nil {
		t.Errorf("nowcoder GetRecentContest failed, err: %v", err)
		return
	}
	t.Logf("nowcoder GetRecentContest result: %s", jsonx.MarshalString(resp))
}
