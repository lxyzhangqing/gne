# gne
本项目基于[《基于最大文本密度的网页正文抽取方法》](https://patents.google.com/patent/CN103714176A/zh)的理论，参考项目[GeneralNewsExtractor](https://github.com/GeneralNewsExtractor/GeneralNewsExtractor)实现。

#### 使用示例

```go
package main

import (
	"fmt"
	"github.com/lxyzhangqing/gne"
	"github.com/lxyzhangqing/gne/extractor"
)

func main() {
	url := "http://xxxxx" // you url
	pageContent := `xxx`  // read the page content with url

	result, err := gne.GeneralNewsExtract(pageContent, &extractor.Options{
		BodyXPath: "//*[@id='article_content']",
		Host:      url,
	})
	
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("标题：", result.Title)
	fmt.Println("作者：", result.Author)
	fmt.Println("发布时间：", result.PublishTime)
	fmt.Println("内容：", result.Content)
	fmt.Println("图片：", result.Images)
	fmt.Println("元数据：", result.Meta)
}
```

