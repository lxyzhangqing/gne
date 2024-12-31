package extractor

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/lxyzhangqing/gne/def"
	"github.com/lxyzhangqing/gne/utils"
	"golang.org/x/net/html"
	"math"
	"strings"
)

func NewContentExtractor(root *html.Node) ContentExtractor {
	return ContentExtractor{
		root:                     root,
		contentTag:               "p",
		punctuation:              `！，。？、；：“”‘’《》%（）,.?:;'"!%()`,
		nodeInfo:                 make(map[string]NodeInfo),
		elementTextCache:         make(map[string][]string),
		highWeightKeywordPattern: def.HighWeightAttrKeyword,
	}
}

type ContentExtractor struct {
	root                     *html.Node
	contentTag               string
	highWeightKeywordPattern []string
	punctuation              string
	nodeInfo                 map[string]NodeInfo
	elementTextCache         map[string][]string
}

func (c ContentExtractor) Extract(opt *Options) (string, []string, error) {
	// 1、网页预处理
	// 1.1 字符编码处理
	// 1.2 网页规范化（标签没有关闭、属性值没有使用单引用或双引号包含、特殊字符没有转义）
	// 2、将网页转换成一棵dom树，并根据特定标签，将网页中的“标签文本块”抽取出来
	// 3、计算最大文本密度
	// 4、抽取正文

	if opt.BodyXPath == "" {
		// 对 HTML 进行预处理可能会破坏 HTML 原有的结构，导致根据原始 HTML 编写的 XPath 不可用
		// 因此，仅对未指定bodyXpath进行预处理，否则直接按给定xpath进行提取
		opt.BodyXPath = "//body"
		utils.RemoveNoiseNode(c.root.FirstChild, opt.NoiseNodeList)
		utils.NormalizeNode(c.root.FirstChild)
	}

	bodyNodes := htmlquery.Find(c.root, opt.BodyXPath)
	if len(bodyNodes) == 0 {
		return "", nil, fmt.Errorf("xpath %s not found", opt.BodyXPath)
	}

	// 遍历节点
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type != html.ElementNode {
			return
		}
		//fmt.Println(">>>", utils.GetNodeRootTree(n))

		nodeHash := utils.GetNodeHash(n)
		densityInfo := c.calcTextDensity(n)
		if densityInfo == nil {
			return
		}

		textDensity := densityInfo.density
		tiText := densityInfo.tiText
		textTagCount := c.countTextTag(n, "p")
		sbdi := c.calcSbDi(tiText, densityInfo.ti, densityInfo.lti)

		imagesList := make([]string, 0)
		imagesNodeList := htmlquery.Find(n, ".//img/@src")
		if imagesNodeList != nil {
			for _, imgNode := range imagesNodeList {
				img := utils.GetNodeText(imgNode)
				imagesList = append(imagesList, img)
			}
		}

		if opt.Host != "" {
			for i := 0; i < len(imagesList); i++ {
				imagesList[i] = utils.PadHostForImages(opt.Host, imagesList[i])
			}
		}

		nodeInfo := NodeInfo{
			ti:           densityInfo.ti,
			lti:          densityInfo.lti,
			tgi:          densityInfo.tgi,
			ltgi:         densityInfo.ltgi,
			node:         n,
			density:      textDensity,
			text:         tiText,
			images:       imagesList,
			textTagCount: textTagCount,
			sbdi:         sbdi,
		}

		c.nodeInfo[nodeHash] = nodeInfo

		// 后计算子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(bodyNodes[0])
	// 计算得分
	c.calcNewScore()

	// 排序，找到分值最高的内容
	var score float64
	var hash string
	for nodeHash, nodeInfo := range c.nodeInfo {
		if nodeInfo.score > score {
			hash = nodeHash
			score = nodeInfo.score
		}
	}
	return c.nodeInfo[hash].text, c.nodeInfo[hash].images, nil
}

func (c ContentExtractor) countTextTag(node *html.Node, tag string) int {
	/**
	当前标签下面的 text()和 p 标签，都应该进行统计
	*/
	if node == nil {
		return 0
	}

	if tag == "" {
		tag = "p"
	}

	var tagNum, directText int

	tagNodes := htmlquery.Find(node, fmt.Sprintf(".//%s", tag))
	if len(tagNodes) > 0 {
		tagNum = len(utils.GetNodeText(tagNodes[0]))
	}

	textNodes := htmlquery.Find(node, "text()")
	if len(textNodes) > 0 {
		directText = len(utils.GetNodeText(textNodes[0]))
	}

	return tagNum + directText
}

func (c ContentExtractor) getAllTextOfElement(nodes []*html.Node) []string {
	textList := make([]string, 0)
	for _, node := range nodes {
		elementFlag := utils.GetNodeRootTree(node)
		// 直接读取缓存的数据，而不是再重复提取一次
		if _, isExist := c.elementTextCache[elementFlag]; isExist {
			textList = append(textList, c.elementTextCache[elementFlag]...)
		} else {
			elementTextList := make([]string, 0)
			textNodes := htmlquery.Find(node, ".//text()")
			if len(textNodes) == 0 {
				continue
			}

			for _, n := range textNodes {
				text := utils.GetNodeText(n)
				text = strings.TrimSpace(text)
				text = strings.Replace(text, "\n", "", -1)
				if text == "" {
					continue
				}
				elementTextList = append(elementTextList, text)
			}
			c.elementTextCache[elementFlag] = elementTextList
			textList = append(textList, elementTextList...)
		}
	}
	return textList
}

func (c ContentExtractor) needSkipLtGi(ti, lti int) bool {
	/**
	有时候，会出现像维基百科一样，在文字里面加a 标签关键词的情况，例如：

	<div>
	我是正文我是正文我是正文<a href="xxx">关键词1</a>我是正文我是正文我是正文我是正文
	我是正文我是正文我是正文我是正文我是正文<a href="xxx">关键词2</a>我是正文我是正文
	我是正文
	</div>

	在这种情况下，tgi = ltgi = 2，计算公式的分母为0. 为了把这种情况和列表页全是链接的
	情况区分出来，所以要做一下判断。检查节点下面所有 a 标签的超链接中的文本数量与本节点
	下面所有文本数量的比值。如果超链接的文本数量占比极少，那么此时，ltgi 应该忽略
	:param ti: 节点 i 的字符串字数
	:param lti: 节点 i 的带链接的字符串字数
	*/
	if lti == 0 {
		return false
	}
	return ti/lti > 10 // 正文的字符数量是链接字符数量的十倍以上
}

func (c ContentExtractor) calcTextDensity(node *html.Node) *TextDensity {
	/**
	根据公式：

		   Ti - LTi
	TDi = -----------
		  TGi - LTGi


	Ti:节点 i 的字符串字数
	LTi：节点 i 的带链接的字符串字数
	TGi：节点 i 的标签数
	LTGi：节点 i 的带连接的标签数
	*/

	if node == nil {
		return nil
	}

	textList := c.getAllTextOfElement([]*html.Node{node})
	if len(textList) == 0 {
		return nil
	}

	tiText := strings.Join(textList, "\n")
	ti := len(tiText)
	ti = c.increaseTagWeight(ti, node)

	aTagList := htmlquery.Find(node, ".//a")
	lti := len(strings.Join(c.getAllTextOfElement(aTagList), ""))

	tgi := len(utils.GetNodeTextByXPath(node, ".//*"))
	ltgi := len(aTagList)

	if tgi-ltgi == 0 {
		if !c.needSkipLtGi(ti, lti) {
			return &TextDensity{
				density: 0,
				tiText:  tiText,
				ti:      ti,
				lti:     lti,
				tgi:     tgi,
				ltgi:    ltgi,
			}
		} else {
			ltgi = 0
		}
	}

	return &TextDensity{
		density: float64(ti-lti) / float64(tgi-ltgi),
		tiText:  tiText,
		ti:      ti,
		lti:     lti,
		tgi:     tgi,
		ltgi:    ltgi,
	}
}

func (c ContentExtractor) increaseTagWeight(ti int, node *html.Node) int {
	class := utils.GetNodeClass(node)
	if class == "" || len(c.highWeightKeywordPattern) == 0 {
		return ti
	}

	for _, pattern := range c.highWeightKeywordPattern {
		if pattern == class {
			return 2 * ti
		}
	}
	return ti
}

func (c ContentExtractor) calcSbDi(text string, ti, lti int) float64 {
	/**
	         Ti - LTi
	SbDi = --------------
			 Sbi + 1

	SbDi: 符号密度
	Sbi：符号数量
	*/
	sbi := c.countPunctuationNum(text)
	sbdi := float64(ti-lti) / float64(sbi+1)
	if sbdi == 0 {
		return 1
	}
	return sbdi
}

func (c ContentExtractor) countPunctuationNum(text string) int {
	count := 0
	for _, char := range text {
		if strings.Contains(c.punctuation, string(char)) {
			count++
		}
	}
	return count
}

func (c ContentExtractor) calcNewScore() {
	/**
	score = 1 * ndi * log10(text_tag_count + 2) * log(sbdi)

	1：在论文里面，这里使用的是 log(std)，但是每一个密度都乘以相同的对数，他们的相对大小是不会改变的，所以我们没有必要计算
	ndi：节点 i 的文本密度
	text_tag_count: 正文所在标签数。例如正文在<p></p>标签里面，这里就是 p 标签数，如果正文在<div></div>标签，这里就是 div 标签数
	sbdi：节点 i 的符号密度
	*/

	for nodeHash, nodeInfo := range c.nodeInfo {
		score := nodeInfo.density * math.Log10(float64(nodeInfo.textTagCount)+2) * math.Log(nodeInfo.sbdi)
		nodeInfo.score = score
		c.nodeInfo[nodeHash] = nodeInfo
	}
}

type TextDensity struct {
	density float64
	tiText  string
	ti      int
	lti     int
	tgi     int
	ltgi    int
}

type NodeInfo struct {
	ti           int // 节点i的字符串字数
	lti          int // 节点i的带链接的字符串字数
	tgi          int // 节点i的标签数
	ltgi         int // 节点i的带链接的标签数
	node         *html.Node
	density      float64
	text         string
	images       []string
	textTagCount int
	sbdi         float64 // 符号密度
	score        float64
}
