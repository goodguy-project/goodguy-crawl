package luogu

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/samber/lo"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/util/httpx"
	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

var (
	regexRecentContest = regexp.MustCompile(`decodeURIComponent\(".*"\)`)
)

func GetRecentContest(_ *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	request, err := http.NewRequest("GET", "https://www.luogu.com.cn/contest/list", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[[]byte]("luogu", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	text, err := func() (string, error) {
		text := regexRecentContest.Find(response)
		text = text[len(`decodeURIComponent("`) : len(text)-len(`")`)]
		return url.QueryUnescape(string(text))
	}()
	if err != nil {
		return nil, errorx.New(err)
	}
	type Host struct {
		Id         int64  `json:"id,omitempty"`
		Name       string `json:"name"`
		IsPremium  bool   `json:"isPremium,omitempty"`
		Uid        int64  `json:"uid,omitempty"`
		Slogan     string `json:"slogan,omitempty"`
		IsAdmin    bool   `json:"isAdmin,omitempty"`
		IsBanned   bool   `json:"isBanned,omitempty"`
		Color      string `json:"color,omitempty"`
		CcfLevel   int64  `json:"ccfLevel,omitempty"`
		Background string `json:"background,omitempty"`
	}
	type Result struct {
		RuleType           int64  `json:"ruleType"`
		VisibilityType     int64  `json:"visibilityType"`
		InvitationCodeType int64  `json:"invitationCodeType"`
		Rated              bool   `json:"rated"`
		Host               *Host  `json:"host"`
		ProblemCount       int64  `json:"problemCount"`
		Id                 int64  `json:"id"`
		Name               string `json:"name"`
		StartTime          int64  `json:"startTime"`
		EndTime            int64  `json:"endTime"`
	}
	type Contests struct {
		Result  []*Result `json:"result"`
		PerPage int64     `json:"perPage"`
		Count   int64     `json:"count"`
	}
	type CurrentData struct {
		Contests *Contests `json:"contests"`
	}
	type Response struct {
		Code            int64        `json:"code"`
		CurrentTemplate string       `json:"currentTemplate"`
		CurrentData     *CurrentData `json:"currentData"`
		CurrentTitle    string       `json:"currentTitle"`
		CurrentTime     int64        `json:"currentTime"`
	}
	data, err := jsonx.Unmarshal[*Response](text)
	if err != nil {
		return nil, errorx.New(err)
	}
	return &proto.GetRecentContestResponse{
		RecentContest: lo.Map(data.CurrentData.Contests.Result,
			func(contest *Result, _ int) *proto.GetRecentContestResponse_Contest {
				return &proto.GetRecentContestResponse_Contest{
					Name:      contest.Name,
					Url:       fmt.Sprintf("https://www.luogu.com.cn/contest/%d", contest.Id),
					Timestamp: contest.StartTime,
					Duration:  int32(contest.EndTime - contest.StartTime),
				}
			}),
		Platform: "luogu",
	}, nil
}
