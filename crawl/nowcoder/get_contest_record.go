package nowcoder

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/util/httpx"
)

func GetContestRecord(req *proto.GetContestRecordRequest) (*proto.GetContestRecordResponse, error) {
	request, err := http.NewRequest("GET",
		fmt.Sprintf("https://ac.nowcoder.com/acm/contest/rating-history?uid=%s", url.QueryEscape(req.GetHandle())),
		bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	type Data struct {
		ContestId   int     `json:"contestId"`
		Rating      float64 `json:"rating"`
		Rank        int     `json:"rank"`
		ChangeValue float64 `json:"changeValue"`
		Time        int64   `json:"time"`
		ContestName string  `json:"contestName"`
		ColorLevel  int     `json:"colorLevel"`
	}
	type Response struct {
		Msg  string  `json:"msg"`
		Code int     `json:"code"`
		Data []*Data `json:"data"`
	}
	response, _, err := httpx.SendRequest[*Response]("nowcoder", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	ret := make([]*proto.GetContestRecordResponse_Record, 0)
	for _, data := range response.Data {
		ret = append(ret, &proto.GetContestRecordResponse_Record{
			Name:      data.ContestName,
			Url:       fmt.Sprintf("https://ac.nowcoder.com/acm/contest/%d", data.ContestId),
			Timestamp: data.Time / 1000,
			Rating:    int32(data.Rating),
		})
	}
	var rating int32 = 0
	if len(ret) > 0 {
		rating = ret[len(ret)-1].Rating
	}
	return &proto.GetContestRecordResponse{
		ProfileUrl: fmt.Sprintf("https://ac.nowcoder.com/acm/home/%s", url.QueryEscape(req.GetHandle())),
		Rating:     rating,
		Length:     int32(len(ret)),
		Record:     ret,
		Platform:   "nowcoder",
		Handle:     req.GetHandle(),
	}, nil
}
