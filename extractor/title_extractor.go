package extractor

import (
	"github.com/antchfx/htmlquery"
	"github.com/lxyzhangqing/gne/def"
	"github.com/lxyzhangqing/gne/utils"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func NewTitleExtractor(root *html.Node) TitleExtractor {
	return TitleExtractor{root: root}
}

type TitleExtractor struct {
	root *html.Node
}

func (t TitleExtractor) extractByXPath(xpath string) string {
	if xpath != "" && t.root != nil {
		titles := htmlquery.Find(t.root, xpath)
		if len(titles) > 0 {
			return utils.GetNodeText(titles[0])
		}
	}
	return ""
}

func (t TitleExtractor) extractByTitle() string {
	titleList := htmlquery.Find(t.root, "//title/text()")
	if len(titleList) == 0 {
		return ""
	}

	re := regexp.MustCompile(def.TitleSplitCharPattern)
	titles := re.Split(titleList[0].Data, -1)
	if len(titles) > 0 {
		if len(titles[0]) > 4 {
			return titles[0]
		}
		return titleList[0].Data
	}
	return ""
}

func (t TitleExtractor) extractByHTag() string {
	titles := htmlquery.Find(t.root, def.TitleHTagXpath)
	if len(titles) > 0 {
		return titles[0].Data
	}
	return ""
}

func (t TitleExtractor) extractByHTagAndTitle() string {
	/**
	一般来说，我们可以认为 title 中包含新闻标题，但是可能也含有其他文字，例如：
	GNE 成为全球最好的新闻提取模块-今日头条
	新华网：GNE 成为全球最好的新闻提取模块

	同时，新闻的某个 <h>标签中也会包含这个新闻标题。

	因此，通过 h 标签与 title 的文字双向匹配，找到最适合作为新闻标题的字符串。
	但是，需要考虑到 title 与 h 标签中的文字可能均含有特殊符号，因此，不能直接通过
	判断 h 标签中的文字是否在 title 中来判断，这里需要中最长公共子串。
	*/
	hTagNodeList := htmlquery.Find(t.root, `(//h1//text() | //h2//text() | //h3//text() | //h4//text() | //h5//text())`)
	titleNodeList := htmlquery.Find(t.root, `//title/text()`)

	if len(hTagNodeList) == 0 || len(titleNodeList) == 0 {
		return ""
	}

	hTagTextList := utils.GetNodesText(hTagNodeList)
	titleTextList := utils.GetNodesText(titleNodeList)
	var newsTitle string

	for _, hTagText := range hTagTextList {
		lcs := utils.GetLongestCommonSubString(titleTextList[0], hTagText)
		if len(lcs) > len(newsTitle) {
			newsTitle = lcs
		}
	}

	if len(newsTitle) > 4 {
		return newsTitle
	}
	return ""
}

func (t TitleExtractor) Extract(xpath string) (string, error) {
	if title := t.extractByXPath(xpath); title != "" {
		return normalize(title), nil
	}

	if title := t.extractByHTagAndTitle(); title != "" {
		return normalize(title), nil
	}

	if title := t.extractByTitle(); title != "" {
		return normalize(title), nil
	}

	if title := t.extractByHTag(); title != "" {
		return normalize(title), nil
	}

	return "", nil
}

func normalize(title string) string {
	title = strings.Replace(title, "\n", "", -1)
	title = strings.Trim(title, " ")
	return title
}
