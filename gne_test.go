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
		NoiseNodeList: []string{`//div[@class="comment-list"]`},
		//PublishTimeXPath: "//*[@id=\"app\"]/div/div[1]/div/div[2]/div[3]/div/div/div/div[1]/div/div[1]/div[1]/div/div/div[1]/div/div[1]/span",
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

func TestContentExtract(t *testing.T) {
	htmlContent := `<body><div><div class="leading-8 text-[#242424] post-content mt-12 text-lg space-y-7"><p>社交媒体巨头Meta正在酝酿一场重大变革。据《金融时报》报道，Meta产品副总裁康纳·海斯透露，公司正在积极开发AI生成角色系统，这些虚拟角色未来将在Facebook等社交平台上扮演重要角色。</p><p>自去年7月推出AI角色创建工具以来，Meta已成功孵化数十万个AI角色。这些虚拟角色目前虽然处于私密状态，但未来将逐步走向社交平台前台。根据海斯的描述，这些AI角色将拥有完整的个人档案，包括个性化头像，并具备内容创作和分享功能，其存在感将不亚于真实用户账号。</p><p style="text-align: center"><img src="https://pic.chinaz.com/picmap/202308101450093085_0.jpg" title="数字人 虚拟主播 (1) (图片来源：AI合成)" alt="数字人 虚拟主播 (1)"></p><p style="text-align: center;">图源备注：图片由AI生成，图片授权服务商Midjourney</p><p>为推进这一技术布局，Meta在未来两年内将把AI技术发展列为核心战略之一。目前，用户已经可以借助Meta AI进行图片编辑或创建AI助手。公司还计划推出文本转视频生成工具，帮助创作者制作AI视频内容。</p><p>然而，这一雄心勃勃的计划也引发了业界担忧。亿万美元男孩公司首席营销官贝基·欧文指出，AI角色缺乏真实的生活经历和情感共鸣，可能会影响平台内容质量。同时，专家们还警告，如果缺乏有效监管，AI账号可能会加剧虚假信息的传播。</p><p>对此，Meta已采取初步防范措施，要求所有AI生成内容必须标注"AI信息"标签。Meta发言人向《福克斯商业》表示，AI工作室将允许用户根据个人兴趣创建并与AI角色互动，这些虚拟角色可用于实用和娱乐等多个场景。</p><p>尽管Meta在AI领域投入力度空前，但其首席执行官马克·扎克伯格也坦言，这项技术的实际收益可能还需要数年时间才能显现。他强调，虽然打造<span class="spamTxt">顶尖</span>AI技术是一项浩大工程，但从长远来看，这笔投资终将为公司和投资者带来可观回报。</p></div></div></body>`
	result, err := GeneralNewsExtract(htmlContent, &extractor.Options{})

	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("内容：", result.Content)
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
