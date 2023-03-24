package codeforces

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/samber/lo"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/util/httpx"
)

func GetRecentContest(_ *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	type Result struct {
		Id                  int64  `json:"id"`
		Name                string `json:"name"`
		Type                string `json:"type"`
		Phase               string `json:"phase"`
		Frozen              bool   `json:"frozen"`
		DurationSeconds     int64  `json:"durationSeconds"`
		Description         string `json:"description,omitempty"`
		Difficulty          int64  `json:"difficulty,omitempty"`
		Kind                string `json:"kind,omitempty"`
		Season              string `json:"season,omitempty"`
		StartTimeSeconds    int64  `json:"startTimeSeconds,omitempty"`
		RelativeTimeSeconds int64  `json:"relativeTimeSeconds,omitempty"`
		PreparedBy          string `json:"preparedBy,omitempty"`
		Country             string `json:"country,omitempty"`
		City                string `json:"city,omitempty"`
		IcpcRegion          string `json:"icpcRegion,omitempty"`
		WebsiteUrl          string `json:"websiteUrl,omitempty"`
	}
	type Response struct {
		Status string    `json:"status"`
		Result []*Result `json:"result"`
	}
	req, err := http.NewRequest("GET", "https://codeforces.com/api/contest.list?gym=false", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	resp, _, err := httpx.SendRequest[*Response]("codeforces", nil, req)
	if err != nil {
		return nil, errorx.New(err)
	}
	return &proto.GetRecentContestResponse{
		RecentContest: lo.Map(resp.Result, func(result *Result, _ int) *proto.GetRecentContestResponse_Contest {
			return &proto.GetRecentContestResponse_Contest{
				Name:      result.Name,
				Url:       fmt.Sprintf("https://codeforces.com/contest/%d", result.Id),
				Timestamp: result.StartTimeSeconds,
				Duration:  int32(result.DurationSeconds),
				ExtInfo:   make(map[string]string),
			}
		}),
		Platform: "codeforces",
	}, nil
}
