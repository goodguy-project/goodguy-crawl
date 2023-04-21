package atcoder

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/antchfx/htmlquery"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/htmlx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
)

func GetContestRecord(req *proto.GetContestRecordRequest) (*proto.GetContestRecordResponse, error) {
	request, err := http.NewRequest("GET",
		fmt.Sprintf("https://atcoder.jp/users/%s/history", url.QueryEscape(req.GetHandle())), bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[[]byte]("atcoder", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	doc, err := htmlquery.Parse(bytes.NewBuffer(response))
	if err != nil {
		return nil, errorx.New(err)
	}
	nodes, err := htmlquery.QueryAll(doc, `//table[@id="history"]//tr[contains(@class, "text-center")]`)
	if err != nil {
		return nil, errorx.New(err)
	}
	ret := make([]*proto.GetContestRecordResponse_Record, 0)
	for _, node := range nodes {
		performance, err := htmlx.GetNodeString(node, `//td[4]`)
		if err != nil {
			return nil, errorx.New(err)
		}
		if performance == "-" {
			continue
		}
		name, err := htmlx.GetNodeString(node, `//td[2]/a[1]`)
		if err != nil {
			return nil, errorx.New(err)
		}
		link, err := htmlx.GetNodeString(node, `//td[2]/a[1]/@href`)
		if err != nil {
			return nil, errorx.New(err)
		}
		ratingStr, err := htmlx.GetNodeString(node, `//td[5]/span`)
		if err != nil {
			return nil, errorx.New(err)
		}
		rating, err := strconv.ParseInt(ratingStr, 10, 64)
		if err != nil {
			return nil, errorx.New(err)
		}
		tsStr, err := htmlx.GetNodeString(node, `//td[1]/@data-order`)
		if err != nil {
			return nil, errorx.New(err)
		}
		ts, err := time.ParseInLocation("2006/01/02 15:04:05", tsStr, location)
		if err != nil {
			return nil, errorx.New(err)
		}
		ret = append(ret, &proto.GetContestRecordResponse_Record{
			Name:      name,
			Url:       "https://atcoder.jp" + link,
			Timestamp: ts.Unix(),
			Rating:    int32(rating),
		})
	}
	var rating int32 = 0
	if len(ret) > 0 {
		rating = ret[len(ret)-1].Rating
	}
	return &proto.GetContestRecordResponse{
		ProfileUrl: "https://atcoder.jp/users/" + url.QueryEscape(req.GetHandle()),
		Rating:     rating,
		Length:     int32(len(ret)),
		Record:     ret,
		Platform:   "atcoder",
		Handle:     req.GetHandle(),
	}, nil
}
