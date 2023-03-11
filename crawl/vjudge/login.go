package vjudge

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/goodguy-project/goodguy-crawl/util/errorx"
	"github.com/goodguy-project/goodguy-crawl/util/httpx"
)

func login(username, password string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errorx.New(err)
	}
	client := &http.Client{Jar: jar}
	request, err := http.NewRequest("POST", "https://vjudge.net/user/login",
		strings.NewReader(fmt.Sprintf("username=%s&password=%s", url.QueryEscape(username), url.QueryEscape(password))))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	if err != nil {
		return nil, errorx.New(err)
	}
	response, _, err := httpx.SendRequest[[]byte]("vjudge", client, request)
	if err != nil {
		return nil, errorx.New(err)
	}
	if r := string(response); r != "success" {
		return nil, errorx.New(errors.New(r))
	}
	return client, nil
}
