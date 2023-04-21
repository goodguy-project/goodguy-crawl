package acwing

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/antchfx/htmlquery"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
)

var (
	regexWeekContest = regexp.MustCompile(`第 \d+ 场周赛`)
	regexStartTime   = regexp.MustCompile(`let start_time = Date.parse\(".*".replace\(" ", "T"\)\);`)
	regexEndTime     = regexp.MustCompile(`let end_time = Date.parse\(".*".replace\(" ", "T"\)\);`)
)

func getContestDuration(url string) (int64, error) {
	request, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return 0, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[[]byte]("acwing", nil, request)
	if err != nil {
		return 0, errorx.New(err)
	}
	startTime := regexStartTime.Find(response)
	if len(startTime) == 0 {
		return 0, errorx.New(nil)
	}
	endTime := regexEndTime.Find(response)
	if len(endTime) == 0 {
		return 0, errorx.New(nil)
	}
	startTime = startTime[len(`let start_time = Date.parse("`) : len(startTime)-len(`".replace(" ", "T"));`)]
	endTime = endTime[len(`let end_time = Date.parse("`) : len(endTime)-len(`".replace(" ", "T"));`)]
	startTime = bytes.ReplaceAll(startTime, []byte(" "), []byte("T"))
	endTime = bytes.ReplaceAll(endTime, []byte(" "), []byte("T"))
	start, err := time.Parse(time.RFC3339Nano, string(startTime))
	if err != nil {
		return 0, errorx.New(err)
	}
	end, err := time.Parse(time.RFC3339Nano, string(endTime))
	if err != nil {
		return 0, errorx.New(err)
	}
	return end.Unix() - start.Unix(), nil
}

func GetRecentContest(_ *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	// 只爬取第一页
	request, err := http.NewRequest("GET", "https://www.acwing.com/activity/1/competition/", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[[]byte]("acwing", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	doc, err := htmlquery.Parse(bytes.NewBuffer(response))
	if err != nil {
		return nil, errorx.New(err)
	}
	contests := make([]*proto.GetRecentContestResponse_Contest, 0)
	nodes, err := htmlquery.QueryAll(doc, `//div[@class="activity-index-block"]`)
	if err != nil {
		return nil, errorx.New(err)
	}
	for _, node := range nodes {
		tsNode, err := htmlquery.QueryAll(node, `//div[@class="row"]//div[@class="col-xs-6"]//span[@class="activity_td"]`)
		if err != nil {
			return nil, errorx.New(err)
		}
		if len(tsNode) <= 1 {
			return nil, errors.New("acwing get contest time failed")
		}
		ts, err := time.ParseInLocation("2006-01-02 15:04:05", htmlquery.OutputHTML(tsNode[1], false), location)
		if err != nil {
			return nil, errorx.New(err)
		}
		nameNode, err := htmlquery.Query(node, `//span[@class="activity_title"]`)
		if err != nil {
			return nil, errorx.New(err)
		}
		urlNode, err := htmlquery.Query(node, `//div[@class="col-md-11"]/a/@href`)
		if err != nil {
			return nil, errorx.New(err)
		}
		name := htmlquery.OutputHTML(nameNode, false)
		url := fmt.Sprintf("https://www.acwing.com/%s", htmlquery.OutputHTML(urlNode, false))
		duration, err := func() (int64, error) {
			if regexWeekContest.MatchString(name) {
				// 周赛持续时间直接返回75分钟
				return 75 * 60, nil
			}
			duration, err := getContestDuration(url)
			if err != nil {
				return 0, errorx.New(err)
			}
			return duration, nil
		}()
		if err != nil {
			return nil, errorx.New(err)
		}
		contests = append(contests, &proto.GetRecentContestResponse_Contest{
			Name:      name,
			Url:       url,
			Timestamp: ts.Unix(),
			Duration:  int32(duration),
		})
	}
	return &proto.GetRecentContestResponse{
		RecentContest: contests,
		Platform:      "acwing",
	}, nil
}
