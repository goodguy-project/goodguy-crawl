package vjudge

import (
	"testing"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func TestGetSubmitData(t *testing.T) {
	resp, err := GetSubmitRecord(&proto.GetSubmitRecordRequest{
		Platform: "vjudge",
		Handle:   "ConanYu",
		AuthInfo: &proto.AuthInfo{
			Username: "",
			Password: "",
		},
	})
	if err != nil {
		t.Errorf("vjudge GetSubmitRecord failed, err: %v", err)
		return
	}
	t.Logf("vjudge GetSubmitRecord result: %s", jsonx.MarshalString(resp))
}
