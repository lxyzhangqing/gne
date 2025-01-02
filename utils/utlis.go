package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/lxyzhangqing/gne/def"
	"golang.org/x/net/html"
	urlpkg "net/url"
	"strings"
)

func GetLongestCommonSubString(str1, str2 string) string {
	if str1 == "" || str2 == "" {
		return ""
	}

	matrix := make([][]int, 0)
	for i := 0; i <= len(str1); i++ {
		matrix = append(matrix, make([]int, len(str2)+1))
	}

	maxLength := 0
	startPosition := 0
	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			if str1[i-1] == str2[j-1] {
				matrix[i][j] = matrix[i-1][j-1] + 1
				if matrix[i][j] > maxLength {
					maxLength = matrix[i][j]
					startPosition = i - maxLength
				}
			} else {
				matrix[i][j] = 0
			}
		}
	}
	return str1[startPosition : startPosition+maxLength]
}

func GetNodeText(node *html.Node) string {
	if node == nil {
		return ""
	}

	var buf bytes.Buffer
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}

	f(node)
	return buf.String()
}

func GetNodeTextByXPath(node *html.Node, xpath string) string {
	if node == nil {
		return ""
	}

	tmpNodes := htmlquery.Find(node, xpath)
	if len(tmpNodes) == 0 {
		return ""
	}

	return GetNodeText(tmpNodes[0])
}

func GetNodesText(nodes []*html.Node) []string {
	if len(nodes) == 0 {
		return []string{}
	}

	texts := make([]string, 0)
	for _, node := range nodes {
		text := GetNodeText(node)
		texts = append(texts, text)
	}
	return texts
}

func GetNodeHtml(node *html.Node) (string, error) {
	if node == nil {
		return "", nil
	}

	var buf bytes.Buffer
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(&buf, c); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func GetNodesTextByXPath(nodes []*html.Node, xpath string) string {
	if len(nodes) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for _, node := range nodes {
		text := GetNodeTextByXPath(node, xpath)
		if text == "" {
			buf.WriteString(text)
		}
	}
	return buf.String()
}

func GetNodeRootTree(node *html.Node) string {
	if node == nil {
		return ""
	}

	var buf bytes.Buffer
	var f func(*html.Node)
	f = func(n *html.Node) {
		// if this node has parent, first to deal parent node
		if n.Parent != nil {
			f(n.Parent)
		}

		// Deal this node
		// Get the index of this node within siblings
		index := 0
		for c := n.PrevSibling; c != nil; c = c.PrevSibling {
			if c.Type == html.TextNode {
				continue
			}

			if c.Data == n.Data {
				index++
			}
		}

		if index != 0 {
			// if the node has the same tag before it, output the index of this node
			buf.WriteString(fmt.Sprintf("/%v[%v]", n.Data, index+1))
		} else {
			isFound := false
			for c := n.NextSibling; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					continue
				}

				if c.Data == n.Data {
					isFound = true
				}
			}

			if isFound {
				// if the node has the same tag behind it, output the index with 1
				buf.WriteString(fmt.Sprintf("/%v[%v]", n.Data, 1))
			} else {
				// if the node no same tag behind it, not to output the index
				buf.WriteString(fmt.Sprintf("/%v", n.Data))
			}
		}
	}
	f(node)
	return buf.String()
}

func GetNodeClass(node *html.Node) string {
	if node == nil || len(node.Attr) == 0 {
		return ""
	}

	for _, attr := range node.Attr {
		if attr.Key == "class" {
			return attr.Val
		}
	}
	return ""
}

func GetNodeClassList(node *html.Node) []string {
	classList := make([]string, 0)
	if node == nil || len(node.Attr) == 0 {
		return classList
	}

	for _, attr := range node.Attr {
		if attr.Key == "class" {
			classList = append(classList, attr.Val)
		}
	}
	return classList
}

func GetNodeHash(node *html.Node) string {
	if node == nil {
		return ""
	}

	tree := GetNodeRootTree(node)
	h := md5.New()
	h.Write([]byte(tree))
	return hex.EncodeToString(h.Sum(nil))
}

func RemoveNoiseNode(node *html.Node, noiseNodeList []string) {
	if node == nil || len(noiseNodeList) == 0 {
		return
	}

	for _, noiseNode := range noiseNodeList {
		nodes := htmlquery.Find(node, noiseNode)
		for _, n := range nodes {
			RemoveNode(n)
		}
	}
}

func NormalizeNode(node *html.Node) {
	StripNodes(node, def.UselessTag)

	do := func(n *html.Node) {
		// remove comment tag
		if n.Type == html.CommentNode {
			RemoveNode(n)
		}

		// inspired by readability.
		if def.TagsCanBeRemoveIfEmpty[n.Data] && n.FirstChild == nil {
			RemoveNode(n)
		}

		// merge text in span or strong to parent p tag
		if n.Data == "p" {
			StripTags(n, "span")
			StripTags(n, "strong")
		}

		// if a div tag does not contain any sub node, it could be converted to p node.
		if n.Data == "div" {
			func() {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type != html.TextNode {
						return
					}
				}
				n.Data = "p"
			}()
		}
		if n.Data == "span" {
			func() {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type != html.TextNode {
						return
					}
				}
				n.Data = "p"
			}()
		}

		// remove empty p tag
		if n.Data == "p" {
			nodes := htmlquery.Find(n, ".//img")
			if len(nodes) == 0 {
				text := GetNodeText(n)
				if text == "" || strings.TrimSpace(text) == "" {
					DropTag(n)
				}
			}
		}

		classList := GetNodeClassList(n)
		for _, class := range classList {
			if def.UseLessAttr[class] {
				RemoveNode(n)
				break
			}
		}
	}

	var fn func(*html.Node)
	fn = func(n *html.Node) {
		do(n)

		// 如果当前节点已经被删除，不遍历其子节点
		if n.Parent == nil && n.PrevSibling == nil && n.NextSibling == nil {
			return
		}

		for c := n.FirstChild; c != nil; {
			next := c.NextSibling
			fn(c)
			c = next
		}
	}
	fn(node)
}

func RemoveNode(node *html.Node) {
	if node.Parent != nil {
		node.Parent.RemoveChild(node)
	}
}

func StripNodes(root *html.Node, tagNamesMap map[string]bool) {
	if root == nil || tagNamesMap == nil {
		return
	}

	var fn func(*html.Node)
	fn = func(n *html.Node) {
		for c := n.FirstChild; c != nil; {
			next := c.NextSibling
			if c.Type != html.TextNode && def.UselessTag[c.Data] {
				n.RemoveChild(c)
			} else {
				fn(c)
			}
			c = next
		}
	}
	fn(root)
}

func DropTag(node *html.Node) {
	/**
	Remove the tag, but not its children or text.  The children and text
	are merged into the parent.
	*/
	if node.Parent != nil {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			node.Parent.InsertBefore(CloneNode(c), node)
		}
		node.Parent.RemoveChild(node)
	}
}

func CloneNode(node *html.Node) *html.Node {
	if node == nil {
		return nil
	}

	return &html.Node{
		Type:       node.Type,
		DataAtom:   node.DataAtom,
		Data:       node.Data,
		Namespace:  node.Namespace,
		Attr:       node.Attr[:],
		FirstChild: node.FirstChild,
		LastChild:  node.LastChild,
	}
}

func StripTags(root *html.Node, tag string) {
	/**
	Delete all elements with the provided tag names from a tree or
	subtree.  This will remove the elements and their attributes, but
	*not* their text/tail content or descendants.  Instead, it will
	merge the text content and children of the element into its parent.
	*/
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n.Type == html.TextNode {
			return
		}

		// child first
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}

		if n.Data == tag {
			if n.Parent != nil {
				// merge the text content and children of the element into its parent
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					n.Parent.InsertBefore(CloneNode(c), n)
				}
				// remove the elements and their attributes
				n.Parent.RemoveChild(n)
			}
		}
	}
	fn(root)
}

func PadHostForImages(host, url string) string {
	/**
	网站上的图片可能有如下几种格式：

	完整的绝对路径：https://xxx.com/1.jpg
	完全不含 host 的相对路径： /1.jpg
	含 host 但是不含 scheme:  xxx.com/1.jpg 或者  ://xxx.com/1.jpg
	*/

	if strings.Index(url, "http") == 0 {
		return url
	}

	parsedUri, err := urlpkg.Parse(host)
	if err != nil {
		return url
	}

	if strings.Index(url, ":") == 0 {
		return fmt.Sprintf("%s%s", parsedUri.Scheme, url)
	}

	if strings.Index(url, "//") == 0 {
		return fmt.Sprintf("%s:%s", parsedUri.Scheme, url)
	}

	result, err := urlpkg.JoinPath(host, url)
	if err != nil {
		return url
	}
	return result
}
