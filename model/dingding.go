package model

type DingDingMarkdown struct {
	MsgType   string    `json:"msgtype"`
	Dmarkdown Dmarkdown `json:"markdown"`
}

type Dmarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func NewDingDingMarkdown(title, text string) *DingDingMarkdown {
	return &DingDingMarkdown{
		MsgType: "markdown",
		Dmarkdown: Dmarkdown{
			Title: title,
			Text:  text,
		},
	}
}
