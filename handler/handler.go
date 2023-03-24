package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goodguy-project/goodguy-crawl/crawl/acwing"
	"github.com/goodguy-project/goodguy-crawl/crawl/atcoder"
	"github.com/goodguy-project/goodguy-crawl/crawl/codechef"
	"github.com/goodguy-project/goodguy-crawl/crawl/codeforces"
	"github.com/goodguy-project/goodguy-crawl/crawl/leetcode"
	"github.com/goodguy-project/goodguy-crawl/crawl/luogu"
	"github.com/goodguy-project/goodguy-crawl/crawl/nowcoder"
	"github.com/goodguy-project/goodguy-crawl/crawl/vjudge"
	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/cachex"
	"github.com/goodguy-project/goodguy-crawl/util/errorx"
)

type (
	getContestRecordFunction func(*proto.GetContestRecordRequest) (*proto.GetContestRecordResponse, error)
	getSubmitRecordFunction  func(*proto.GetSubmitRecordRequest) (*proto.GetSubmitRecordResponse, error)
	getRecentContestFunction func(*proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error)
	getDailyQuestionFunction func(*proto.GetDailyQuestionRequest) (*proto.GetDailyQuestionResponse, error)
)

var (
	getContestRecordMap = map[string]getContestRecordFunction{
		"codeforces": cachex.TTLWrap(codeforces.GetContestRecord, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"atcoder":    cachex.TTLWrap(atcoder.GetContestRecord, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"nowcoder":   cachex.TTLWrap(nowcoder.GetContestRecord, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"leetcode":   cachex.TTLWrap(leetcode.GetContestRecord, cachex.TTLConfig{TTL: 2 * time.Hour}),
	}
	getSubmitRecordMap = map[string]getSubmitRecordFunction{
		"codeforces": cachex.TTLWrap(codeforces.GetSubmitRecord, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"luogu":      cachex.TTLWrap(luogu.GetSubmitRecord, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"vjudge":     cachex.TTLWrap(vjudge.GetSubmitRecord, cachex.TTLConfig{TTL: 2 * time.Hour}),
	}
	getRecentContestMap = map[string]getRecentContestFunction{
		"codeforces": cachex.TTLWrap(codeforces.GetRecentContest, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"atcoder":    cachex.TTLWrap(atcoder.GetRecentContest, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"nowcoder":   cachex.TTLWrap(nowcoder.GetRecentContest, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"luogu":      cachex.TTLWrap(luogu.GetRecentContest, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"leetcode":   cachex.TTLWrap(leetcode.GetRecentContest, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"codechef":   cachex.TTLWrap(codechef.GetRecentContest, cachex.TTLConfig{TTL: 2 * time.Hour}),
		"acwing":     cachex.TTLWrap(acwing.GetRecentContest, cachex.TTLConfig{TTL: 2 * time.Hour}),
	}
	getDailyQuestionMap = map[string]getDailyQuestionFunction{
		"leetcode": cachex.TTLWrap(leetcode.GetDailyQuestion, cachex.TTLConfig{TTL: 2 * time.Hour}),
	}
)

func GetContestRecord(_ context.Context, req *proto.GetContestRecordRequest) (*proto.GetContestRecordResponse, error) {
	function, ok := getContestRecordMap[req.GetPlatform()]
	if !ok {
		return nil, errorx.New(errors.New(fmt.Sprintf("GetContestRecord of platform (%s) not implemented", req.GetPlatform())))
	}
	return function(req)
}

func GetSubmitRecord(_ context.Context, req *proto.GetSubmitRecordRequest) (*proto.GetSubmitRecordResponse, error) {
	function, ok := getSubmitRecordMap[req.GetPlatform()]
	if !ok {
		return nil, errorx.New(errors.New(fmt.Sprintf("GetSubmitRecord of platform (%s) not implemented", req.GetPlatform())))
	}
	return function(req)
}

func GetRecentContest(_ context.Context, req *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	function, ok := getRecentContestMap[req.GetPlatform()]
	if !ok {
		return nil, errorx.New(errors.New(fmt.Sprintf("GetRecentContest of platform (%s) not implemented", req.GetPlatform())))
	}
	return function(req)
}

func GetDailyQuestion(_ context.Context, req *proto.GetDailyQuestionRequest) (*proto.GetDailyQuestionResponse, error) {
	function, ok := getDailyQuestionMap[req.GetPlatform()]
	if !ok {
		return nil, errorx.New(errors.New(fmt.Sprintf("GetDailyQuestion of platform (%s) not implemented", req.GetPlatform())))
	}
	return function(req)
}
