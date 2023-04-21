package leetcode

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/samber/lo"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/jsonx"
)

func GetRecentContest(_ *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errorx.New(err)
	}
	client := &http.Client{Jar: jar}
	csrfToken := getCsrfToken(client, "https://leetcode.cn/contest/")
	request, err := http.NewRequest("POST", "https://leetcode.cn/graphql",
		bytes.NewBuffer(jsonx.Marshal(map[string]interface{}{
			"operationName": nil,
			"variables":     make(map[string]interface{}),
			"query":         "{\n  contestUpcomingContests {\n    containsPremium\n    title\n    cardImg\n    titleSlug\n    description\n    startTime\n    duration\n    originStartTime\n    isVirtual\n    isLightCardFontColor\n    company {\n      watermark\n      __typename\n    }\n    __typename\n  }\n}\n",
		})))
	if err != nil {
		return nil, errorx.New(err)
	}
	request.Header.Set("x-csrftoken", csrfToken)
	request.Header.Set("origin", "https://leetcode.cn")
	request.Header.Set("referer", "https://leetcode.cn/contest/")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json")

	type ContestUpcomingContests struct {
		ContainsPremium      bool   `json:"containsPremium"`
		Title                string `json:"title"`
		CardImg              string `json:"cardImg"`
		TitleSlug            string `json:"titleSlug"`
		Description          string `json:"description"`
		StartTime            int64  `json:"startTime"`
		Duration             int64  `json:"duration"`
		OriginStartTime      int64  `json:"originStartTime"`
		IsVirtual            bool   `json:"isVirtual"`
		IsLightCardFontColor bool   `json:"isLightCardFontColor"`
		Typename             string `json:"__typename"`
	}
	type Data struct {
		ContestUpcomingContests []*ContestUpcomingContests `json:"contestUpcomingContests"`
	}
	type Request struct {
		Data *Data `json:"data"`
	}
	response, _, err := httpx.SendRequest[*Request]("leetcode", client, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	return &proto.GetRecentContestResponse{
		RecentContest: lo.Map(response.Data.ContestUpcomingContests, func(contest *ContestUpcomingContests, _ int) *proto.GetRecentContestResponse_Contest {
			return &proto.GetRecentContestResponse_Contest{
				Name:      contest.Title,
				Url:       fmt.Sprintf("https://leetcode.cn/contest/%s", contest.TitleSlug),
				Timestamp: contest.StartTime,
				Duration:  int32(contest.Duration),
			}
		}),
		Platform: "leetcode",
	}, nil
}
