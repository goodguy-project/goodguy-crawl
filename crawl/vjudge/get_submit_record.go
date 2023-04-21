package vjudge

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/tidwall/gjson"

	"github.com/goodguy-project/goodguy-crawl/v2/proto"
	"github.com/goodguy-project/goodguy-crawl/v2/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
)

var (
	verdictMap = map[string]proto.GetSubmitRecordResponse_SubmitRecord_Verdict{
		"AC":  proto.GetSubmitRecordResponse_SubmitRecord_Accepted,
		"WA":  proto.GetSubmitRecordResponse_SubmitRecord_WrongAnswer,
		"CE":  proto.GetSubmitRecordResponse_SubmitRecord_CompilationError,
		"TLE": proto.GetSubmitRecordResponse_SubmitRecord_TimeLimitExceeded,
		"MLE": proto.GetSubmitRecordResponse_SubmitRecord_MemoryLimitExceeded,
		"RE":  proto.GetSubmitRecordResponse_SubmitRecord_RuntimeError,
		"OLE": proto.GetSubmitRecordResponse_SubmitRecord_OutputLimitExceeded,
		"PE":  proto.GetSubmitRecordResponse_SubmitRecord_PresentationError,
	}
	programmingLanguageMap = map[string]proto.GetSubmitRecordResponse_SubmitRecord_ProgrammingLanguage{
		"CPP":    proto.GetSubmitRecordResponse_SubmitRecord_Cpp,
		"C":      proto.GetSubmitRecordResponse_SubmitRecord_C,
		"JAVA":   proto.GetSubmitRecordResponse_SubmitRecord_Java,
		"PYTHON": proto.GetSubmitRecordResponse_SubmitRecord_Python,
		"RUBY":   proto.GetSubmitRecordResponse_SubmitRecord_Ruby,
		"CSHARP": proto.GetSubmitRecordResponse_SubmitRecord_CSharp,
	}
)

func GetSubmitRecord(req *proto.GetSubmitRecordRequest) (*proto.GetSubmitRecordResponse, error) {
	username, password := func() (string, string) {
		if username, password := req.GetAuthInfo().GetUsername(), req.GetAuthInfo().GetPassword(); username != "" && password != "" {
			return username, password
		}
		if username, password := os.Getenv("VJUDGE_USERNAME"), os.Getenv("VJUDGE_PASSWORD"); username != "" && password != "" {
			return username, password
		}
		return "", ""
	}()
	if username == "" && password == "" {
		return nil, errorx.New(errors.New("vjudge GetSubmitRecord need auth info"))
	}
	client, err := login(req.GetAuthInfo().GetUsername(), req.GetAuthInfo().GetPassword())
	if err != nil {
		return nil, errorx.New(err)
	}
	problemSet := mapset.NewSet[string]()
	maxId := ""
	submitCount := 0
	ojDistribution := make(map[string]int32)
	submitRecordData := make([]*proto.GetSubmitRecordResponse_SubmitRecord, 0)
	for true {
		request, err := http.NewRequest("GET",
			fmt.Sprintf("https://vjudge.net/user/submissions?username=%s&pageSize=500&maxId=%s",
				url.QueryEscape(req.GetHandle()), url.QueryEscape(maxId)),
			bytes.NewBuffer([]byte{}))
		if err != nil {
			return nil, errorx.New(err)
		}
		response, _, err := httpx.SendRequest[[]byte]("vjudge", client, request)
		if err != nil {
			return nil, errorx.New(err)
		}
		array := gjson.Get(string(response), "data").Array()
		if len(array) == 0 {
			break
		}
		for _, submit := range array {
			data := submit.Array()
			if len(data) <= 9 {
				return nil, errorx.New(nil)
			}
			pb := fmt.Sprintf("%s-%s", data[2].String(), data[3].String())
			if !problemSet.Contains(pb) && data[4].String() == "AC" {
				problemSet.Add(pb)
				ojDistribution[data[2].String()] = ojDistribution[data[2].String()] + 1
			}
			submitRecordData = append(submitRecordData, &proto.GetSubmitRecordResponse_SubmitRecord{
				ProblemName: pb,
				ProblemUrl:  fmt.Sprintf("https://vjudge.net/problem/%s", pb),
				SubmitTime:  data[9].Int() / 1000,
				RunningTime: int32(data[5].Int()),
				Verdict: func() proto.GetSubmitRecordResponse_SubmitRecord_Verdict {
					if verdict, ok := verdictMap[data[4].String()]; ok {
						return verdict
					}
					return proto.GetSubmitRecordResponse_SubmitRecord_Other
				}(),
				ProgrammingLanguage: func() proto.GetSubmitRecordResponse_SubmitRecord_ProgrammingLanguage {
					if programmingLanguage, ok := programmingLanguageMap[data[7].String()]; ok {
						return programmingLanguage
					}
					return proto.GetSubmitRecordResponse_SubmitRecord_Unknown
				}(),
			})
		}
		submitCount += len(array)
		maxId = strconv.FormatInt(array[len(array)-1].Array()[0].Int()-1, 10)
	}
	return &proto.GetSubmitRecordResponse{
		ProfileUrl:     fmt.Sprintf("https://vjudge.net/user/%s", req.GetHandle()),
		AcceptCount:    int32(problemSet.Cardinality()),
		SubmitCount:    int32(submitCount),
		Distribution:   nil,
		OjDistribution: ojDistribution,
		Platform:       "vjudge",
		Handle:         req.GetHandle(),
		SubmitRecord:   submitRecordData,
	}, nil
}
