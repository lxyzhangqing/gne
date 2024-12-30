package extractor

import (
	"github.com/antchfx/htmlquery"
	"github.com/lxyzhangqing/gne/utils"
	"golang.org/x/net/html"
)

func NewMetaExtractor(root *html.Node) MetaExtractor {
	return MetaExtractor{root: root}
}

type MetaExtractor struct {
	root *html.Node
}

func (m MetaExtractor) Extract() (map[string]string, error) {
	metaContent := make(map[string]string)

	metaNodes := htmlquery.Find(m.root, "//meta")
	if metaNodes == nil {
		return metaContent, nil
	}

	for _, metaNode := range metaNodes {
		var name, content string
		nameNodes := htmlquery.Find(metaNode, "@name|@property")
		if len(nameNodes) == 0 {
			continue
		}
		name = utils.GetNodeText(nameNodes[0])

		contentNodes := htmlquery.Find(metaNode, "@content")
		if len(contentNodes) == 0 {
			continue
		}
		content = utils.GetNodeText(contentNodes[0])
		metaContent[name] = content
	}

	return metaContent, nil
}
