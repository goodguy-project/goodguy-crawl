package atcoder

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/jsonx"
)

func TestGetContestRecord(t *testing.T) {
	resp, err := GetContestRecord(&proto.GetContestRecordRequest{
		Platform: "atcoder",
		Handle:   "conanyu",
	})
	if err != nil {
		t.Errorf("atcoder GetContestRecord failed, err: %v", err)
		return
	}
	t.Logf("atcoder GetContestRecord result: %s", jsonx.MarshalString(resp))
}
