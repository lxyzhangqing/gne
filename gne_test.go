package gne

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/lxyzhangqing/gne/extractor"
	"os"
	"strings"
	"testing"
)

func readHtml() (string, error) {
	fileName := "./test.html"
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func TestExtract(t *testing.T) {
	htmlContent, err := readHtml()
	if err != nil {
		t.Fatal(err)
	}

	result, err := GeneralNewsExtract(htmlContent, &extractor.Options{
		PublishTimeXPath: "//*[@id=\"app\"]/div/div[1]/div/div[2]/div[3]/div/div/div/div[1]/div/div[1]/div[1]/div/div/div[1]/div/div[1]/span",
	})
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("标题：", result.Title)
	fmt.Println("作者：", result.Author)
	fmt.Println("发布时间：", result.PublishTime)
	fmt.Println("内容：", result.Content)
	fmt.Println("图片：", result.Images)
	fmt.Println("元数据：", result.Meta)
}

func TestListPageExtract(t *testing.T) {
	htmlContent, err := readHtml()
	if err != nil {
		t.Fatal(err)
	}

	list, err := ListPageExtract(htmlContent)
	for _, l := range list {
		fmt.Println(l)
	}
}

func TestXpath(t *testing.T) {
	htmlContent, err := readHtml()
	if err != nil {
		t.Fatal(err)
	}

	top, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatal(err)
	}

	nodes := htmlquery.Find(top, "//title/text()")
	if len(nodes) > 0 {
		fmt.Printf("nodes: %v\n", len(nodes))
		fmt.Println("TITLE: ", nodes[0].Data)
	}
}
