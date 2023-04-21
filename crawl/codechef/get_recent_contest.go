package codechef

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
)

func GetRecentContest(_ *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	params := url.Values{}
	params.Set("sort_by", "START")
	params.Set("sorting_order", "asc")
	params.Set("offset", "0")
	params.Set("mode", "premium")
	request, err := http.NewRequest("GET", "https://www.codechef.com/api/list/contests/all?"+params.Encode(), bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	type Contest struct {
		ContestCode         string    `json:"contest_code"`
		ContestName         string    `json:"contest_name"`
		ContestStartDate    string    `json:"contest_start_date"`
		ContestEndDate      string    `json:"contest_end_date"`
		ContestStartDateIso time.Time `json:"contest_start_date_iso"`
		ContestEndDateIso   time.Time `json:"contest_end_date_iso"`
		ContestDuration     string    `json:"contest_duration"`
		DistinctUsers       int       `json:"distinct_users"`
	}
	type Banners struct {
		Image string `json:"image"`
		Link  string `json:"link"`
	}
	type Response struct {
		Status           string     `json:"status"`
		Message          string     `json:"message"`
		PresentContests  []*Contest `json:"present_contests"`
		FutureContests   []*Contest `json:"future_contests"`
		PracticeContests []*Contest `json:"practice_contests"`
		PastContests     []*Contest `json:"past_contests"`
		Banners          []*Banners `json:"banners"`
	}
	response, _, err := httpx.SendRequest[*Response]("codechef", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	ret := make([]*proto.GetRecentContestResponse_Contest, 0)
	trans := func(contest *Contest) *proto.GetRecentContestResponse_Contest {
		duration, _ := strconv.ParseInt(contest.ContestDuration, 10, 64)
		return &proto.GetRecentContestResponse_Contest{
			Name:      contest.ContestCode,
			Url:       fmt.Sprintf("https://www.codechef.com/%s", contest.ContestCode),
			Timestamp: contest.ContestStartDateIso.Unix(),
			ExtInfo:   nil,
			Duration:  int32(duration * 60),
		}
	}
	for _, contest := range response.PresentContests {
		ret = append(ret, trans(contest))
	}
	for _, contest := range response.FutureContests {
		ret = append(ret, trans(contest))
	}
	for _, contest := range response.PracticeContests {
		ret = append(ret, trans(contest))
	}
	for _, contest := range response.PastContests {
		ret = append(ret, trans(contest))
	}
	return &proto.GetRecentContestResponse{
		RecentContest: ret,
		Platform:      "codechef",
	}, nil
}
