package htmlx

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"

	"github.com/goodguy-project/goodguy-crawl/util/errorx"
)

func GetNodeString(node *html.Node, xpath string) (string, error) {
	query, err := htmlquery.Query(node, xpath)
	if query == nil || err != nil {
		return "", errorx.New(err)
	}
	return htmlquery.OutputHTML(query, false), nil
}
