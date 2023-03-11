package codeforces

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func TestGetContestRecord(t *testing.T) {
	resp, err := GetContestRecord(&proto.GetContestRecordRequest{
		Platform: "codeforces",
		Handle:   "conanyu",
	})
	if err != nil {
		t.Errorf("codeforces GetContestRecord failed, err: %v", err)
		return
	}
	t.Logf("codeforces GetContestRecord result: %s", jsonx.MarshalString(resp))
}
