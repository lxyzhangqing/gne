package extractor

type Options struct {
	TitleXPath       string
	AuthorXPath      string
	PublishTimeXPath string
	BodyXPath        string
	Normalize        bool
	NoiseNodeList    []string
	WithBodyHtml     bool
	UseVisibleInfo   bool
	Host             string
}
