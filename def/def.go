package def

var AuthorPattern = []string{
	"责编[：|:| |丨|/]\\s*([\u4E00-\u9FA5a-zA-Z]{2,20})",
	"责任编辑[：|:| |丨|/]\\s*([\u4E00-\u9FA5a-zA-Z]{2,20})",
	"作者[：|:| |丨|/]\\s*([\u4E00-\u9FA5a-zA-Z]{2,20})",
	"编辑[：|:| |丨|/]\\s*([\u4E00-\u9FA5a-zA-Z]{2,20})",
	"文[：|:| |丨|/]\\s*([\u4E00-\u9FA5a-zA-Z]{2,20})",
	"原创[：|:| |丨|/]\\s*([\u4E00-\u9FA5a-zA-Z]{2,20})",
	"撰文[：|:| |丨|/]\\s*([\u4E00-\u9FA5a-zA-Z]{2,20})",
	"来源[：|:| |丨|/]\\s*([\u4E00-\u9FA5a-zA-Z]{2,20})",
	// 以下正则表达式需要进一步测试
	// '(作者[：|:| |丨|/]\s*[\u4E00-\u9FA5a-zA-Z、 ]{2,20})[）】)]]?[^\u4E00-\u9FA5|:|：]',
	//'(记者[：|:| |丨|/]\s*[\u4E00-\u9FA5a-zA-Z、 ]{2,20})[）】)]]?[^\u4E00-\u9FA5|:|：]',
	//'(原创[：|:| |丨|/]\s*[\u4E00-\u9FA5a-zA-Z、 ]{2,20})[）】)]]?[^\u4E00-\u9FA5|:|：]',
	// '(撰文[：|:| |丨|/]\s*[\u4E00-\u9FA5a-zA-Z、 ]{2,20})[）】)]]?[^\u4E00-\u9FA5|:|：]',
	// '(文/图[：|:| |丨|/]?\s*[\u4E00-\u9FA5a-zA-Z、 ]{2,20})[）】)]]?[^\u4E00-\u9FA5|:|：]',
}

var DatetimePattern = []string{
	`(\d{4}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[0-1]?[0-9]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{4}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[2][0-3]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{4}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[0-1]?[0-9]:[0-5]?[0-9])`,
	`(\d{4}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[2][0-3]:[0-5]?[0-9])`,
	`(\d{4}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[1-24]\d时[0-60]\d分)([1-24]\d时)`,
	`(\d{2}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[0-1]?[0-9]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{2}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[2][0-3]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{2}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[0-1]?[0-9]:[0-5]?[0-9])`,
	`(\d{2}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[2][0-3]:[0-5]?[0-9])`,
	`(\d{2}[-|/|.]\d{1,2}[-|/|.]\d{1,2}\s*?[1-24]\d时[0-60]\d分)([1-24]\d时)`,
	`(\d{4}年\d{1,2}月\d{1,2}日\s*?[0-1]?[0-9]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{4}年\d{1,2}月\d{1,2}日\s*?[2][0-3]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{4}年\d{1,2}月\d{1,2}日\s*?[0-1]?[0-9]:[0-5]?[0-9])`,
	`(\d{4}年\d{1,2}月\d{1,2}日\s*?[2][0-3]:[0-5]?[0-9])`,
	`(\d{4}年\d{1,2}月\d{1,2}日\s*?[1-24]\d时[0-60]\d分)([1-24]\d时)`,
	`(\d{2}年\d{1,2}月\d{1,2}日\s*?[0-1]?[0-9]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{2}年\d{1,2}月\d{1,2}日\s*?[2][0-3]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{2}年\d{1,2}月\d{1,2}日\s*?[0-1]?[0-9]:[0-5]?[0-9])`,
	`(\d{2}年\d{1,2}月\d{1,2}日\s*?[2][0-3]:[0-5]?[0-9])`,
	`(\d{2}年\d{1,2}月\d{1,2}日\s*?[1-24]\d时[0-60]\d分)([1-24]\d时)`,
	`(\d{1,2}月\d{1,2}日\s*?[0-1]?[0-9]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{1,2}月\d{1,2}日\s*?[2][0-3]:[0-5]?[0-9]:[0-5]?[0-9])`,
	`(\d{1,2}月\d{1,2}日\s*?[0-1]?[0-9]:[0-5]?[0-9])`,
	`(\d{1,2}月\d{1,2}日\s*?[2][0-3]:[0-5]?[0-9])`,
	`(\d{1,2}月\d{1,2}日\s*?[1-24]\d时[0-60]\d分)([1-24]\d时)`,
	`(\d{4}[-|/|.]\d{1,2}[-|/|.]\d{1,2})`,
	`(\d{2}[-|/|.]\d{1,2}[-|/|.]\d{1,2})`,
	`(\d{4}年\d{1,2}月\d{1,2}日)`,
	`(\d{2}年\d{1,2}月\d{1,2}日)`,
	`(\d{1,2}月\d{1,2}日)`,
}

const TitleHTagXpath = `//h1//text() | //h2//text() | //h3//text() | //h4//text()`

const TitleSplitCharPattern = `[-_|]`

var UselessTag = map[string]bool{
	"style":      true,
	"script":     true,
	"link":       true,
	"video":      true,
	"iframe":     true,
	"source":     true,
	"picture":    true,
	"header":     true,
	"blockquote": true,
	"footer":     true,
}

// TagsCanBeRemoveIfEmpty
// if one tag in the follow list does not contain any child node nor content, it could be removed
var TagsCanBeRemoveIfEmpty = map[string]bool{
	"section": true,
	"h1":      true,
	"h2":      true,
	"h3":      true,
	"h4":      true,
	"h5":      true,
	"h6":      true,
	"span":    true,
}

var UseLessAttr = map[string]bool{
	"share":          true,
	"contribution":   true,
	"copyright":      true,
	"copy-right":     true,
	"disclaimer":     true,
	"recommend":      true,
	"related":        true,
	"footer":         true,
	"comment":        true,
	"social":         true,
	"submeta":        true,
	"report-infor":   true,
	"header_toolbar": true,
}

var HighWeightAttrKeyword = []string{"content", "article", "news_txt", "pages_content", "post_text"}

// PublishTimeMeta 部分特别规范的新闻网站，可以直接从 HTML 的 meta 数据中获得发布时间
var PublishTimeMeta = []string{
	`//meta[starts-with(@property, "rnews:datePublished")]/@content`,
	`//meta[starts-with(@property, "article:published_time")]/@content`,
	`//meta[starts-with(@property, "og:published_time")]/@content`,
	`//meta[starts-with(@property, "og:release_date")]/@content`,
	`//meta[starts-with(@itemprop, "datePublished")]/@content`,
	`//meta[starts-with(@itemprop, "dateUpdate")]/@content`,
	`//meta[starts-with(@name, "OriginalPublicationDate")]/@content`,
	`//meta[starts-with(@name, "article_date_original")]/@content`,
	`//meta[starts-with(@name, "og:time")]/@content`,
	`//meta[starts-with(@name, "apub:time")]/@content`,
	`//meta[starts-with(@name, "publication_date")]/@content`,
	`//meta[starts-with(@name, "sailthru.date")]/@content`,
	`//meta[starts-with(@name, "PublishDate")]/@content`,
	`//meta[starts-with(@name, "publishdate")]/@content`,
	`//meta[starts-with(@name, "PubDate")]/@content`,
	`//meta[starts-with(@name, "pubtime")]/@content`,
	`//meta[starts-with(@name, "_pubtime")]/@content`,
	`//meta[starts-with(@name, "weibo: article:create_at")]/@content`,
	`//meta[starts-with(@pubdate, "pubdate")]/@content`,
}

// ArticleXpath 满足下面的XPath，极有可能是文章详情页
var ArticleXpath = []string{`//*[@class="article__content"]`}
