package leetcode

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/samber/lo"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/jsonx"
)

func GetContestRecord(req *proto.GetContestRecordRequest) (*proto.GetContestRecordResponse, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errorx.New(err)
	}
	client := &http.Client{Jar: jar}
	csrfToken := getCsrfToken(client, fmt.Sprintf("https://leetcode.cn/u/%s/", url.QueryEscape(req.GetHandle())))
	request, err := http.NewRequest("POST", "https://leetcode.cn/graphql",
		bytes.NewBuffer(jsonx.Marshal(map[string]interface{}{
			"operationName": "userContest",
			"variables": map[string]string{
				"userSlug": req.GetHandle(),
			},
			"query": "query userContest($userSlug: String!) {\n  userContestRanking(userSlug: $userSlug) {\n    currentRatingRanking\n    ratingHistory\n    levelHistory\n    contestRankingHistoryV2\n    contestHistory\n    __typename\n  }\n  globalRatingRanking(userSlug: $userSlug)\n  userContestScore(userSlug: $userSlug)\n  contestUnratedContests\n}\n",
		})))
	if err != nil {
		return nil, errorx.New(err)
	}
	request.Header.Set("x-csrftoken", csrfToken)
	request.Header.Set("x-definition-name", "userContestRanking,globalRatingRanking,userContestScore,contestUnratedContests")
	request.Header.Set("x-operation-name", "userContest")
	request.Header.Set("x-timezone", "Asia/Shanghai")
	request.Header.Set("Content-Type", "application/json")
	type UserContestRanking struct {
		CurrentRatingRanking    int    `json:"currentRatingRanking"`
		RatingHistory           string `json:"ratingHistory"`
		LevelHistory            string `json:"levelHistory"`
		ContestRankingHistoryV2 string `json:"contestRankingHistoryV2"`
		ContestHistory          string `json:"contestHistory"`
		Typename                string `json:"__typename"`
	}
	type Data struct {
		UserContestRanking     *UserContestRanking `json:"userContestRanking"`
		GlobalRatingRanking    int                 `json:"globalRatingRanking"`
		UserContestScore       string              `json:"userContestScore"`
		ContestUnratedContests string              `json:"contestUnratedContests"`
	}
	type Response struct {
		Data *Data `json:"data"`
	}
	response, _, err := httpx.SendRequest[*Response]("leetcode", client, request)
	if response == nil || err != nil {
		return nil, errorx.New(err)
	}
	type ContestTitle struct {
		Title     string `json:"title"`
		TitleSlug string `json:"title_slug"`
	}
	if response.Data == nil || response.Data.UserContestRanking == nil {
		return nil, errorx.New(nil)
	}
	contestTitle, err := jsonx.Unmarshal[[]*ContestTitle](response.Data.UserContestRanking.ContestHistory)
	if err != nil {
		return nil, errorx.New(err)
	}
	contestHistory := make(map[string]*lo.Tuple2[int, string])
	lo.ForEach(contestTitle, func(item *ContestTitle, index int) {
		contestHistory[item.TitleSlug] = &lo.Tuple2[int, string]{A: index, B: item.Title}
	})
	ratingHistory, err := jsonx.Unmarshal[[]*float64](response.Data.UserContestRanking.RatingHistory)
	if err != nil {
		return nil, errorx.New(err)
	}
	type UserContestScore struct {
		TitleSlug   string `json:"title_slug"`
		StartTime   int    `json:"start_time"`
		Score       int    `json:"score"`
		FinishTime  int    `json:"finish_time"`
		TimePenalty int    `json:"time_penalty"`
	}
	userContestScore, err := jsonx.Unmarshal[[]*UserContestScore](response.Data.UserContestScore)
	if err != nil {
		return nil, errorx.New(err)
	}
	ret := make([]*proto.GetContestRecordResponse_Record, 0)
	for _, r := range userContestScore {
		contest := contestHistory[r.TitleSlug]
		rating := ratingHistory[contest.A]
		ret = append(ret, &proto.GetContestRecordResponse_Record{
			Name:      contest.B,
			Url:       fmt.Sprintf("https://leetcode.cn/contest/%s", r.TitleSlug),
			Timestamp: int64(r.StartTime),
			Rating: func() int32 {
				if rating == nil {
					return 1500
				}
				return int32(*rating)
			}(),
		})
	}
	var rating int32 = 1500
	if len(ratingHistory) > 0 {
		if r := ratingHistory[len(ratingHistory)-1]; r != nil {
			rating = int32(*r)
		}
	}
	return &proto.GetContestRecordResponse{
		ProfileUrl: fmt.Sprintf("https://leetcode.cn/u/%s/", url.QueryEscape(req.GetHandle())),
		Rating:     rating,
		Length:     int32(len(userContestScore)),
		Record:     ret,
		Platform:   "leetcode",
		Handle:     req.GetHandle(),
	}, nil
}
