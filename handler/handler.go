package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/goodguy-project/goodguy-crawl/crawl/acwing"
	"github.com/goodguy-project/goodguy-crawl/crawl/atcoder"
	"github.com/goodguy-project/goodguy-crawl/crawl/codechef"
	"github.com/goodguy-project/goodguy-crawl/crawl/codeforces"
	"github.com/goodguy-project/goodguy-crawl/crawl/leetcode"
	"github.com/goodguy-project/goodguy-crawl/crawl/luogu"
	"github.com/goodguy-project/goodguy-crawl/crawl/nowcoder"
	"github.com/goodguy-project/goodguy-crawl/crawl/vjudge"
	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/cache"
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
		"codeforces": cache.TTLWrap(codeforces.GetContestRecord, cache.TTLConfig{TTL: 7200}),
		"atcoder":    cache.TTLWrap(atcoder.GetContestRecord, cache.TTLConfig{TTL: 7200}),
		"nowcoder":   cache.TTLWrap(nowcoder.GetContestRecord, cache.TTLConfig{TTL: 7200}),
		"leetcode":   cache.TTLWrap(leetcode.GetContestRecord, cache.TTLConfig{TTL: 7200}),
	}
	getSubmitRecordMap = map[string]getSubmitRecordFunction{
		"codeforces": cache.TTLWrap(codeforces.GetSubmitRecord, cache.TTLConfig{TTL: 7200}),
		"luogu":      cache.TTLWrap(luogu.GetSubmitRecord, cache.TTLConfig{TTL: 7200}),
		"vjudge":     cache.TTLWrap(vjudge.GetSubmitRecord, cache.TTLConfig{TTL: 7200}),
	}
	getRecentContestMap = map[string]getRecentContestFunction{
		"codeforces": cache.TTLWrap(codeforces.GetRecentContest, cache.TTLConfig{TTL: 7200}),
		"atcoder":    cache.TTLWrap(atcoder.GetRecentContest, cache.TTLConfig{TTL: 7200}),
		"nowcoder":   cache.TTLWrap(nowcoder.GetRecentContest, cache.TTLConfig{TTL: 7200}),
		"luogu":      cache.TTLWrap(luogu.GetRecentContest, cache.TTLConfig{TTL: 7200}),
		"leetcode":   cache.TTLWrap(leetcode.GetRecentContest, cache.TTLConfig{TTL: 7200}),
		"codechef":   cache.TTLWrap(codechef.GetRecentContest, cache.TTLConfig{TTL: 7200}),
		"acwing":     cache.TTLWrap(acwing.GetRecentContest, cache.TTLConfig{TTL: 7200}),
	}
	getDailyQuestionMap = map[string]getDailyQuestionFunction{
		"leetcode": cache.TTLWrap(leetcode.GetDailyQuestion, cache.TTLConfig{TTL: 7200}),
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
