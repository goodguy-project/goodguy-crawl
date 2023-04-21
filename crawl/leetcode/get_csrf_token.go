package leetcode

import (
	"bytes"
	"net/http"

	"github.com/goodguy-project/goodguy-crawl/v2/util/httpx"
)

func getCsrfToken(client *http.Client, url string) string {
	if client == nil {
		panic("getCsrfToken.client == nil")
	}
	request, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return ""
	}
	_, extra, err := httpx.SendRequest[struct{}]("leetcode", client, request)
	if err != nil {
		return ""
	}
	for _, cookie := range extra.Cookies {
		if cookie.Name == "csrftoken" {
			return cookie.Value
		}
	}
	return ""
}
