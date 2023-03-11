package leetcode

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/util/httpx"
	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func GetDailyQuestion(_ *proto.GetDailyQuestionRequest) (*proto.GetDailyQuestionResponse, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errorx.New(err)
	}
	client := &http.Client{Jar: jar}
	csrfToken := getCsrfToken(client, "https://leetcode.cn/contest/")
	request, err := http.NewRequest("POST", "https://leetcode.cn/graphql",
		bytes.NewBuffer(jsonx.Marshal(map[string]interface{}{
			"query":     "query questionOfToday{todayRecord{date question{questionId difficulty title titleCn: translatedTitle titleSlug status}}}",
			"variables": make(map[string]string),
		})))
	if err != nil {
		return nil, errorx.New(err)
	}
	request.Header.Set("x-csrftoken", csrfToken)
	request.Header.Set("origin", "https://leetcode.cn")
	request.Header.Set("referer", "https://leetcode.cn/contest/")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json")
	type Question struct {
		QuestionID string `json:"questionId"`
		Difficulty string `json:"difficulty"`
		Title      string `json:"title"`
		TitleCn    string `json:"titleCn"`
		TitleSlug  string `json:"titleSlug"`
	}
	type TodayRecord struct {
		Date     string    `json:"date"`
		Question *Question `json:"question"`
	}
	type Data struct {
		TodayRecord []*TodayRecord `json:"todayRecord"`
	}
	type Response struct {
		Data *Data `json:"data"`
	}
	response, _, err := httpx.SendRequest[*Response]("leetcode", client, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	problem := make([]*proto.GetDailyQuestionResponse_Problem, 0)
	for _, r := range response.Data.TodayRecord {
		q := r.Question
		problem = append(problem, &proto.GetDailyQuestionResponse_Problem{
			Platform:   "leetcode",
			Url:        fmt.Sprintf("https://leetcode.cn/problems/%s", q.TitleSlug),
			Id:         q.QuestionID,
			Name:       q.TitleCn,
			Difficulty: q.Difficulty,
		})
	}
	return &proto.GetDailyQuestionResponse{Problem: problem}, nil
}
