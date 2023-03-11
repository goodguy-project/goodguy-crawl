package nowcoder

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func TestGetContestRecord(t *testing.T) {
	resp, err := GetContestRecord(&proto.GetContestRecordRequest{
		Platform: "nowcoder",
		Handle:   "6693394",
	})
	if err != nil {
		t.Errorf("nowcoder GetContestRecord failed, err: %v", err)
		return
	}
	t.Logf("nowcoder GetContestRecord result: %s", jsonx.MarshalString(resp))
}
