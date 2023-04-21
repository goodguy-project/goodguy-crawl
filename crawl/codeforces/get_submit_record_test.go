package codeforces

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/jsonx"
)

func TestGetSubmitData(t *testing.T) {
	resp, err := GetSubmitRecord(&proto.GetSubmitRecordRequest{
		Platform: "codeforces",
		Handle:   "conanyu",
	})
	if err != nil {
		t.Errorf("codeforces GetSubmitRecord failed, err: %v", err)
		return
	}
	t.Logf("codeforces GetSubmitRecord result: %s", jsonx.MarshalString(resp))
}
