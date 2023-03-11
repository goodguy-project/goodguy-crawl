package handler

import (
	"context"
	"testing"

	"github.com/goodguy-project/goodguy-crawl/util/jsonx"

	"github.com/goodguy-project/goodguy-crawl/proto"
)

func TestGetContestRecord(t *testing.T) {
	r, err := GetContestRecord(context.Background(), &proto.GetContestRecordRequest{
		Platform: "codeforces",
		Handle:   "conanyu",
	})
	if err != nil {
		t.Error("GetContestRecord failed")
		return
	}
	t.Logf("GetContestRecord response: %s", jsonx.MarshalString(r))
}

func TestGetDailyQuestion(t *testing.T) {
	r, err := GetDailyQuestion(context.Background(), &proto.GetDailyQuestionRequest{
		Platform: "leetcode",
	})
	if err != nil {
		t.Error("GetDailyQuestion failed")
		return
	}
	t.Logf("GetDailyQuestion response: %s", jsonx.MarshalString(r))
}

func TestGetRecentContest(t *testing.T) {
	r, err := GetRecentContest(context.Background(), &proto.GetRecentContestRequest{
		Platform: "leetcode",
	})
	if err != nil {
		t.Error("GetRecentContest failed")
		return
	}
	t.Logf("GetRecentContest response: %s", jsonx.MarshalString(r))
}

func TestGetSubmitRecord(t *testing.T) {
	r, err := GetSubmitRecord(context.Background(), &proto.GetSubmitRecordRequest{
		Platform: "codeforces",
		Handle:   "conanyu",
	})
	if err != nil {
		t.Error("GetSubmitRecord failed")
		return
	}
	t.Logf("GetSubmitRecord response: %s", jsonx.MarshalString(r))
}
