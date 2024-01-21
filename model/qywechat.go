package model

type QyWeChatMarkdown struct {
	MsgType  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

func NewQyWeChatMarkdown(content string) *QyWeChatMarkdown {
	return &QyWeChatMarkdown{
		MsgType: "markdown",
		Markdown: Markdown{
			Content: content,
		},
	}
}
