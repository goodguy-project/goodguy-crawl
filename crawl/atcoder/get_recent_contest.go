package atcoder

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/antchfx/htmlquery"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/util/httpx"
)

func getContestDuration(url string) int32 {
	request, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return 0
	}
	response, _, err := httpx.SendRequest[[]byte]("atcoder", nil, request)
	if err != nil {
		return 0
	}
	doc, err := htmlquery.Parse(bytes.NewBuffer(response))
	if err != nil {
		return 0
	}
	node, err := htmlquery.Query(doc, `//div[@id="contest-statement"]//span[@class="lang-en"]//ul[1]//li[1]`)
	if err != nil || node == nil {
		return 0
	}
	ret := 0
	for _, c := range htmlquery.OutputHTML(node, false) {
		if '0' <= c && c <= '9' {
			ret = ret*10 + int(c-'0')
		}
	}
	return int32(ret)
}

func GetRecentContest(_ *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	request, err := http.NewRequest("GET", "https://atcoder.jp/?lang=en", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	html, _, err := httpx.SendRequest[[]byte]("atcoder", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	doc, err := htmlquery.Parse(bytes.NewBuffer(html))
	if err != nil {
		return nil, errorx.New(err)
	}
	startTimes, err := htmlquery.QueryAll(doc, `//div[@id="contest-table-upcoming"]//tbody//a[@target="blank"]/time`)
	if err != nil {
		return nil, errorx.New(err)
	}
	contests, err := htmlquery.QueryAll(doc, `//div[@id="contest-table-upcoming"]//tbody//a[name(@target)!="target"]`)
	if err != nil {
		return nil, errorx.New(err)
	}
	urls, err := htmlquery.QueryAll(doc, `//div[@id="contest-table-upcoming"]//tbody//a[name(@target)!="target"]/@href`)
	if err != nil {
		return nil, errorx.New(err)
	}
	length := len(startTimes)
	ret := make([]*proto.GetRecentContestResponse_Contest, 0)
	for idx := 0; idx < length; idx++ {
		startTime, err := time.Parse("2006-01-02 15:04:05-0700", htmlquery.OutputHTML(startTimes[idx], false))
		if err != nil {
			return nil, err
		}
		url := fmt.Sprintf("https://atcoder.jp%s", htmlquery.OutputHTML(urls[idx], false))
		ret = append(ret, &proto.GetRecentContestResponse_Contest{
			Name:      htmlquery.OutputHTML(contests[idx], false),
			Url:       url,
			Timestamp: startTime.Unix(),
			Duration:  getContestDuration(url),
		})
	}
	return &proto.GetRecentContestResponse{
		RecentContest: ret,
		Platform:      "atcoder",
	}, nil
}
