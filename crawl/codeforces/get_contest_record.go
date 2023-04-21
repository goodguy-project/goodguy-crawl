package codeforces

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/samber/lo"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
)

func GetContestRecord(req *proto.GetContestRecordRequest) (*proto.GetContestRecordResponse, error) {
	type Result struct {
		ContestId               int64  `json:"contestId"`
		ContestName             string `json:"contestName"`
		Handle                  string `json:"handle"`
		Rank                    int64  `json:"rank"`
		RatingUpdateTimeSeconds int64  `json:"ratingUpdateTimeSeconds"`
		OldRating               int64  `json:"oldRating"`
		NewRating               int64  `json:"newRating"`
	}
	type Response struct {
		Status string    `json:"status"`
		Result []*Result `json:"result"`
	}
	request, err := http.NewRequest("GET",
		fmt.Sprintf("https://codeforces.com/api/user.rating?handle=%s", url.QueryEscape(req.GetHandle())),
		bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[*Response]("codeforces", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	rating := int64(0)
	if n := len(response.Result); n > 0 {
		rating = response.Result[n-1].NewRating
	}
	return &proto.GetContestRecordResponse{
		ProfileUrl: fmt.Sprintf("https://codeforces.com/profile/%s", url.QueryEscape(req.GetHandle())),
		Rating:     int32(rating),
		Length:     int32(len(response.Result)),
		Record: lo.Map(response.Result, func(contest *Result, _ int) *proto.GetContestRecordResponse_Record {
			return &proto.GetContestRecordResponse_Record{
				Name:      contest.ContestName,
				Url:       fmt.Sprintf("https://codeforces.com/contest/%d", contest.ContestId),
				Timestamp: contest.RatingUpdateTimeSeconds,
				Rating:    int32(contest.NewRating),
			}
		}),
		Platform: "codeforces",
		Handle:   req.GetHandle(),
	}, nil
}
