package luogu

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/util/httpx"
	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func getUserId(handle string) (int, error) {
	request, err := http.NewRequest("GET",
		fmt.Sprintf("https://www.luogu.com.cn/api/user/search?keyword=%s", url.QueryEscape(handle)),
		bytes.NewBuffer([]byte{}))
	if err != nil {
		return 0, errorx.New(err)
	}
	type User struct {
		Uid        int         `json:"uid"`
		Name       string      `json:"name"`
		Slogan     string      `json:"slogan"`
		Badge      interface{} `json:"badge"`
		IsAdmin    bool        `json:"isAdmin"`
		IsBanned   bool        `json:"isBanned"`
		Color      string      `json:"color"`
		CcfLevel   int         `json:"ccfLevel"`
		Background string      `json:"background"`
	}
	type Response struct {
		Users []*User `json:"users"`
	}
	response, _, err := httpx.SendRequest[*Response]("luogu", nil, request)
	if err != nil {
		return 0, errorx.New(err)
	}
	if len(response.Users) == 0 {
		return 0, errorx.New(errors.New("no such user"))
	}
	return response.Users[0].Uid, nil
}

var (
	regexSubmitRecord = regexp.MustCompile(`decodeURIComponent\(.*"\)`)
)

func GetSubmitRecord(req *proto.GetSubmitRecordRequest) (*proto.GetSubmitRecordResponse, error) {
	userId, err := getUserId(req.GetHandle())
	if err != nil {
		return nil, errorx.New(err)
	}
	request, err := http.NewRequest("GET",
		fmt.Sprintf("https://www.luogu.com.cn/user/%d", userId), bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[[]byte]("luogu", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	text, err := func() (string, error) {
		text := regexSubmitRecord.Find(response)
		text = text[len(`decodeURIComponent("`) : len(text)-len(`")`)]
		return url.QueryUnescape(string(text))
	}()
	if err != nil {
		return nil, errorx.New(err)
	}
	type User struct {
		RegisterTime          int    `json:"registerTime,omitempty"`
		Introduction          string `json:"introduction,omitempty"`
		Prize                 []any  `json:"prize,omitempty"`
		FollowingCount        int    `json:"followingCount,omitempty"`
		FollowerCount         int    `json:"followerCount,omitempty"`
		Ranking               any    `json:"ranking,omitempty"`
		BlogAddress           any    `json:"blogAddress,omitempty"`
		PassedProblemCount    int    `json:"passedProblemCount,omitempty"`
		SubmittedProblemCount int    `json:"submittedProblemCount,omitempty"`
		UID                   int    `json:"uid,omitempty"`
		Name                  string `json:"name,omitempty"`
		Slogan                string `json:"slogan,omitempty"`
		Badge                 any    `json:"badge,omitempty"`
		IsAdmin               bool   `json:"isAdmin,omitempty"`
		IsBanned              bool   `json:"isBanned,omitempty"`
		Color                 string `json:"color,omitempty"`
		CcfLevel              int    `json:"ccfLevel,omitempty"`
		Background            string `json:"background,omitempty"`
	}
	type PassedProblems struct {
		Pid        string `json:"pid,omitempty"`
		Title      string `json:"title,omitempty"`
		Difficulty int    `json:"difficulty,omitempty"`
		FullScore  int    `json:"fullScore,omitempty"`
		Type       string `json:"type,omitempty"`
	}
	type SubmittedProblems struct {
		Pid        string `json:"pid,omitempty"`
		Title      string `json:"title,omitempty"`
		Difficulty int    `json:"difficulty,omitempty"`
		FullScore  int    `json:"fullScore,omitempty"`
		Type       string `json:"type,omitempty"`
	}
	type CurrentData struct {
		User              *User                `json:"user,omitempty"`
		PassedProblems    []*PassedProblems    `json:"passedProblems,omitempty"`
		SubmittedProblems []*SubmittedProblems `json:"submittedProblems,omitempty"`
	}
	type Response struct {
		Code            int          `json:"code,omitempty"`
		CurrentTemplate string       `json:"currentTemplate,omitempty"`
		CurrentData     *CurrentData `json:"currentData,omitempty"`
		CurrentTitle    string       `json:"currentTitle,omitempty"`
		CurrentTheme    any          `json:"currentTheme,omitempty"`
		CurrentTime     int          `json:"currentTime,omitempty"`
	}
	msg, err := jsonx.Unmarshal[*Response](text)
	if err != nil {
		return nil, errorx.New(err)
	}
	distribution := make(map[int32]int32)
	acceptProblemSet := mapset.NewSet[string]()
	for _, acceptProblem := range msg.CurrentData.PassedProblems {
		if !acceptProblemSet.Contains(acceptProblem.Pid) {
			acceptProblemSet.Add(acceptProblem.Pid)
			diff := int32(acceptProblem.Difficulty*100 + 100)
			distribution[diff] = distribution[diff] + 1
		}
	}
	return &proto.GetSubmitRecordResponse{
		ProfileUrl:   fmt.Sprintf("https://www.luogu.com.cn/user/%d", userId),
		AcceptCount:  int32(msg.CurrentData.User.PassedProblemCount),
		SubmitCount:  int32(msg.CurrentData.User.SubmittedProblemCount),
		Distribution: distribution,
		Platform:     "luogu",
		Handle:       req.GetHandle(),
		SubmitRecord: nil,
	}, nil
}
