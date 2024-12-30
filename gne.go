package gne

import (
	"github.com/lxyzhangqing/gne/extractor"
	"golang.org/x/net/html"
	"strings"
)

type Result struct {
	Title       string
	Author      string
	PublishTime string
	Content     string
	Images      []string
	Meta        map[string]string
}

func GeneralNewsExtract(htmlContent string, opt *extractor.Options) (*Result, error) {
	root, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	author, err := extractor.NewAuthorExtractor(root).Extract(opt.AuthorXPath)
	if err != nil {
		return nil, err
	}

	meta, err := extractor.NewMetaExtractor(root).Extract()
	if err != nil {
		return nil, err
	}

	publishTime, err := extractor.NewTimeExtractor(root).Extract(opt.PublishTimeXPath)
	if err != nil {
		return nil, err
	}

	title, err := extractor.NewTitleExtractor(root).Extract(opt.TitleXPath)
	if err != nil {
		return nil, err
	}

	content, images, err := extractor.NewContentExtractor(root).Extract(opt)
	if err != nil {
		return nil, err
	}

	return &Result{
		Title:       title,
		Author:      author,
		PublishTime: publishTime,
		Content:     content,
		Images:      images,
		Meta:        meta,
	}, nil
}

func ListPageExtract(htmlContent string) ([]string, error) {
	root, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	return extractor.NewListExtractor(root).Extract()
}
