package extractor

import (
	"github.com/antchfx/htmlquery"
	"github.com/lxyzhangqing/gne/def"
	"github.com/lxyzhangqing/gne/utils"
	"golang.org/x/net/html"
	"regexp"
)

func NewTimeExtractor(root *html.Node) TimeExtractor {
	return TimeExtractor{root: root}
}

type TimeExtractor struct {
	root *html.Node
}

func (t TimeExtractor) Extract(xpath string) (string, error) {
	if pt := t.extractFromUserXpath(xpath); pt != "" {
		return pt, nil
	}

	if pt := t.extractFromMeta(); pt != "" {
		return pt, nil
	}

	if pt := t.extractFromText(); pt != "" {
		return pt, nil
	}

	return "", nil
}

func (t TimeExtractor) extractFromUserXpath(xpath string) string {
	if xpath != "" {
		nodes := htmlquery.Find(t.root, xpath)
		if len(nodes) > 0 {
			return utils.GetNodeText(nodes[0])
		}
	}
	return ""
}

func (t TimeExtractor) extractFromText() string {
	text := utils.GetNodeText(t.root)
	for _, pattern := range def.DatetimePattern {
		re := regexp.MustCompile(pattern)
		result := re.FindString(text)
		if result != "" {
			return result
		}
	}
	return ""
}

func (t TimeExtractor) extractFromMeta() string {
	// 一些很规范的新闻网站，会把新闻的发布时间放在 META 中，因此应该优先检查 META 数据
	for _, xpath := range def.PublishTimeMeta {
		publishTimeNodes := htmlquery.Find(t.root, xpath)
		if len(publishTimeNodes) > 0 {
			return utils.GetNodeText(publishTimeNodes[0])
		}
	}
	return ""
}
