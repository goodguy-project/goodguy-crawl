package nowcoder

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/util/httpx"
)

var (
	timezone, _ = time.LoadLocation("Asia/Shanghai")
)

func getContest(node *html.Node, isOfficial bool) (*proto.GetRecentContestResponse_Contest, error) {
	timeNode, err := htmlquery.Query(node, `//li[@class="match-time-icon"]`)
	if err != nil {
		return nil, errorx.New(err)
	}
	name, err := htmlquery.Query(node, `//a`)
	if err != nil {
		return nil, errorx.New(err)
	}
	href, err := htmlquery.Query(node, `//a/@href`)
	if err != nil {
		return nil, errorx.New(err)
	}
	user, err := htmlquery.Query(node, `//li[@class="user-icon"]`)
	if err != nil {
		return nil, errorx.New(err)
	}
	const layout = "2006-01-02 15:04"
	const length = len(layout)
	getTimeStampFromStr := func(s string) (int64, error) {
		ts, err := time.ParseInLocation(layout, s, location)
		if err != nil {
			return 0, err
		}
		ts = ts.In(timezone)
		return ts.Unix(), nil
	}
	timeText := strings.ReplaceAll(htmlquery.OutputHTML(timeNode, false), "\n", "")
	var unix []int64
	for from := 0; from+length < len(timeText); from++ {
		to := from + length
		ts, err := getTimeStampFromStr(timeText[from:to])
		if err == nil {
			unix = append(unix, ts)
		}
	}
	if len(unix) != 2 {
		return nil, errors.New("nowcoder get contest time failed")
	}
	nameText := htmlquery.OutputHTML(name, false)
	hrefText := htmlquery.OutputHTML(href, false)
	userText := htmlquery.OutputHTML(user, false)
	return &proto.GetRecentContestResponse_Contest{
		Name:      nameText,
		Url:       fmt.Sprintf("https://ac.nowcoder.com%s", hrefText),
		Timestamp: unix[0],
		ExtInfo: map[string]string{
			"user": userText,
			"type": func() string {
				if isOfficial {
					return "official"
				}
				return "unofficial"
			}(),
		},
		Duration: int32(unix[1] - unix[0]),
	}, nil
}

func getUnofficialContest() ([]*proto.GetRecentContestResponse_Contest, error) {
	request, err := http.NewRequest("GET",
		"https://ac.nowcoder.com/acm/contest/vip-index?topCategoryFilter=13", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[[]byte]("nowcoder", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	doc, err := htmlquery.Parse(bytes.NewBuffer(response))
	if err != nil {
		return nil, errorx.New(err)
	}
	contests, err := htmlquery.QueryAll(doc, `//div[contains(@class, "js-current")]//div[@class="platform-item-cont"]`)
	if err != nil {
		return nil, errorx.New(err)
	}
	var ret []*proto.GetRecentContestResponse_Contest
	for _, contest := range contests {
		c, err := getContest(contest, false)
		if err != nil {
			return nil, errorx.New(err)
		}
		ret = append(ret, c)
	}
	return ret, nil
}

func getOfficialContest() ([]*proto.GetRecentContestResponse_Contest, error) {
	request, err := http.NewRequest("GET",
		"https://ac.nowcoder.com/acm/contest/vip-index?topCategoryFilter=14", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[[]byte]("nowcoder", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	doc, err := htmlquery.Parse(bytes.NewBuffer(response))
	if err != nil {
		return nil, errorx.New(err)
	}
	contests, err := htmlquery.QueryAll(doc, `//div[contains(@class, "js-current")]//div[@class="platform-item-cont"]`)
	if err != nil {
		return nil, errorx.New(err)
	}
	var ret []*proto.GetRecentContestResponse_Contest
	for _, contest := range contests {
		c, err := getContest(contest, true)
		if err != nil {
			return nil, errorx.New(err)
		}
		ret = append(ret, c)
	}
	return ret, nil
}

func GetRecentContest(_ *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	contest := make([]*proto.GetRecentContestResponse_Contest, 0)
	unofficialContest, err := getUnofficialContest()
	if err != nil {
		return nil, errorx.New(err)
	}
	officialContest, err := getOfficialContest()
	if err != nil {
		return nil, errorx.New(err)
	}
	contest = append(contest, unofficialContest...)
	contest = append(contest, officialContest...)
	return &proto.GetRecentContestResponse{
		RecentContest: contest,
		Platform:      "nowcoder",
	}, nil
}
