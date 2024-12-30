package extractor

import (
	"golang.org/x/net/html"
)

type ListExtractor struct {
	root *html.Node
}

func (l ListExtractor) Extract() ([]string, error) {
	return nil, nil
}

func NewListExtractor(root *html.Node) ListExtractor {
	return ListExtractor{root: root}
}
