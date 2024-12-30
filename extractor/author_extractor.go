package extractor

import (
	"github.com/antchfx/htmlquery"
	"github.com/lxyzhangqing/gne/def"
	"github.com/lxyzhangqing/gne/utils"
	"golang.org/x/net/html"
	"regexp"
)

type AuthorExtractor struct {
	root *html.Node
}

func NewAuthorExtractor(root *html.Node) AuthorExtractor {
	return AuthorExtractor{root: root}
}

func (a AuthorExtractor) Extract(xpath string) (string, error) {
	if xpath == "" {
		xpath = ".//text()"
	}

	textNodes := htmlquery.Find(a.root, xpath)
	if len(textNodes) == 0 {
		return "", nil
	}

	text := utils.GetNodeText(textNodes[0])
	for _, pattern := range def.AuthorPattern {
		re := regexp.MustCompile(pattern)
		group := re.FindAllStringSubmatch(text, -1)
		if len(group) > 0 && len(group[0]) > 1 {
			return group[0][1], nil
		}
	}
	return "", nil
}
