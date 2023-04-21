package luogu

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/jsonx"
)

func TestGetSubmitData(t *testing.T) {
	resp, err := GetSubmitRecord(&proto.GetSubmitRecordRequest{
		Platform: "luogu",
		Handle:   "yuzining",
	})
	if err != nil {
		t.Errorf("luogu GetSubmitRecord failed, err: %v", err)
		return
	}
	t.Logf("luogu GetSubmitRecord result: %s", jsonx.MarshalString(resp))
}
