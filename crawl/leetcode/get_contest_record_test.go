package leetcode

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func TestGetContestRecord(t *testing.T) {
	resp, err := GetContestRecord(&proto.GetContestRecordRequest{
		Platform: "leetcode",
		Handle:   "johnkram",
	})
	if err != nil {
		t.Errorf("leetcode GetContestRecord failed, err: %v", err)
		return
	}
	t.Logf("leetcode GetContestRecord result: %s", jsonx.MarshalString(resp))
}
