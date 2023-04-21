package codeforces

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
)

var (
	regexC      = regexp.MustCompile(`^GNU GCC`)
	regexCpp    = regexp.MustCompile(`G\+\+|Clang\+\+|C\+\+`)
	regexPython = regexp.MustCompile(`(Pypy|Python)`)
	regexJs     = regexp.MustCompile(`^(JavaScript V8|Node.js)`)
	regexJava   = regexp.MustCompile(`^Java `)
	regexKotlin = regexp.MustCompile(`^Kotlin `)
)

func getProgrammingLanguage(lang string) proto.GetSubmitRecordResponse_SubmitRecord_ProgrammingLanguage {
	if regexC.MatchString(lang) {
		return proto.GetSubmitRecordResponse_SubmitRecord_C
	}
	if regexCpp.MatchString(lang) {
		return proto.GetSubmitRecordResponse_SubmitRecord_Cpp
	}
	if regexPython.MatchString(lang) {
		return proto.GetSubmitRecordResponse_SubmitRecord_Python
	}
	if regexJs.MatchString(lang) {
		return proto.GetSubmitRecordResponse_SubmitRecord_JavaScript
	}
	if regexJava.MatchString(lang) {
		return proto.GetSubmitRecordResponse_SubmitRecord_Java
	}
	if regexKotlin.MatchString(lang) {
		return proto.GetSubmitRecordResponse_SubmitRecord_Kotlin
	}
	if strings.Contains(lang, "C#") {
		return proto.GetSubmitRecordResponse_SubmitRecord_CSharp
	}
	if strings.Contains(lang, "PHP") {
		return proto.GetSubmitRecordResponse_SubmitRecord_PHP
	}
	if strings.Contains(lang, "Ruby") {
		return proto.GetSubmitRecordResponse_SubmitRecord_Ruby
	}
	if strings.Contains(lang, "Scala") {
		return proto.GetSubmitRecordResponse_SubmitRecord_Scala
	}
	if strings.Contains(lang, "Haskell") {
		return proto.GetSubmitRecordResponse_SubmitRecord_Haskell
	}
	return proto.GetSubmitRecordResponse_SubmitRecord_Unknown
}

var (
	verdictMap = map[string]proto.GetSubmitRecordResponse_SubmitRecord_Verdict{
		"Other":                   proto.GetSubmitRecordResponse_SubmitRecord_Other,
		"Accepted":                proto.GetSubmitRecordResponse_SubmitRecord_Accepted,
		"WrongAnswer":             proto.GetSubmitRecordResponse_SubmitRecord_WrongAnswer,
		"RuntimeError":            proto.GetSubmitRecordResponse_SubmitRecord_RuntimeError,
		"TimeLimitExceeded":       proto.GetSubmitRecordResponse_SubmitRecord_TimeLimitExceeded,
		"MemoryLimitExceeded":     proto.GetSubmitRecordResponse_SubmitRecord_MemoryLimitExceeded,
		"CompilationError":        proto.GetSubmitRecordResponse_SubmitRecord_CompilationError,
		"PresentationError":       proto.GetSubmitRecordResponse_SubmitRecord_PresentationError,
		"IdlenessLimitExceeded":   proto.GetSubmitRecordResponse_SubmitRecord_IdlenessLimitExceeded,
		"SecurityViolated":        proto.GetSubmitRecordResponse_SubmitRecord_SecurityViolated,
		"Crashed":                 proto.GetSubmitRecordResponse_SubmitRecord_Crashed,
		"InputPreparationCrashed": proto.GetSubmitRecordResponse_SubmitRecord_InputPreparationCrashed,
		"Partial":                 proto.GetSubmitRecordResponse_SubmitRecord_Partial,
		"Challenged":              proto.GetSubmitRecordResponse_SubmitRecord_Challenged,
		"Skipped":                 proto.GetSubmitRecordResponse_SubmitRecord_Skipped,
		"Testing":                 proto.GetSubmitRecordResponse_SubmitRecord_Testing,
		"Rejected":                proto.GetSubmitRecordResponse_SubmitRecord_Rejected,
		"OutputLimitExceeded":     proto.GetSubmitRecordResponse_SubmitRecord_OutputLimitExceeded,
	}
)

func GetSubmitRecord(req *proto.GetSubmitRecordRequest) (*proto.GetSubmitRecordResponse, error) {
	type Problem struct {
		ContestId int64    `json:"contestId"`
		Index     string   `json:"index"`
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Rating    int64    `json:"rating,omitempty"`
		Tags      []string `json:"tags"`
		Points    float64  `json:"points,omitempty"`
	}
	type Members struct {
		Handle string `json:"handle"`
	}
	type Author struct {
		ContestId        int64      `json:"contestId"`
		Members          []*Members `json:"members"`
		ParticipantType  string     `json:"participantType"`
		Ghost            bool       `json:"ghost"`
		StartTimeSeconds int64      `json:"startTimeSeconds,omitempty"`
		Room             int64      `json:"room,omitempty"`
		TeamId           int64      `json:"teamId,omitempty"`
		TeamName         string     `json:"teamName,omitempty"`
	}
	type Result struct {
		Id                  int64    `json:"id"`
		ContestId           int64    `json:"contestId"`
		CreationTimeSeconds int64    `json:"creationTimeSeconds"`
		RelativeTimeSeconds int64    `json:"relativeTimeSeconds"`
		Problem             *Problem `json:"problem"`
		Author              *Author  `json:"author"`
		ProgrammingLanguage string   `json:"programmingLanguage"`
		Verdict             string   `json:"verdict"`
		Testset             string   `json:"testset"`
		PassedTestCount     int64    `json:"passedTestCount"`
		TimeConsumedMillis  int64    `json:"timeConsumedMillis"`
		MemoryConsumedBytes int64    `json:"memoryConsumedBytes"`
		Points              float64  `json:"points,omitempty"`
	}
	type Response struct {
		Status string    `json:"status"`
		Result []*Result `json:"result"`
	}
	request, err := http.NewRequest("GET",
		fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s", url.QueryEscape(req.GetHandle())),
		bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[*Response]("codeforces", nil, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	acProblemSet := mapset.NewSet[string]()
	distribution := make(map[int32]int32)
	submitRecord := make([]*proto.GetSubmitRecordResponse_SubmitRecord, 0)
	for _, submit := range response.Result {
		problem := fmt.Sprintf("%d-%s", submit.Problem.ContestId, submit.Problem.Index)
		problemRating := submit.Problem.Rating
		submitRecord = append(submitRecord, &proto.GetSubmitRecordResponse_SubmitRecord{
			ProblemName:         submit.Problem.Name,
			ProblemUrl:          fmt.Sprintf("https://codeforces.com/contest/%d/problem/%s", submit.Problem.ContestId, submit.Problem.Index),
			SubmitTime:          submit.CreationTimeSeconds,
			Verdict:             verdictMap[submit.Verdict],
			RunningTime:         int32(submit.TimeConsumedMillis),
			ProgrammingLanguage: getProgrammingLanguage(submit.ProgrammingLanguage),
		})
		if submit.Verdict == "OK" && !acProblemSet.Contains(problem) {
			acProblemSet.Add(problem)
			if problemRating > 0 {
				distribution[int32(problemRating)] = distribution[int32(problemRating)] + 1
			}
		}
	}
	return &proto.GetSubmitRecordResponse{
		ProfileUrl:   fmt.Sprintf("https://codeforces.com/profile/%s", url.QueryEscape(req.GetHandle())),
		AcceptCount:  int32(acProblemSet.Cardinality()),
		SubmitCount:  int32(len(response.Result)),
		Distribution: distribution,
		Platform:     "codeforces",
		Handle:       req.GetHandle(),
		SubmitRecord: submitRecord,
	}, nil
}
