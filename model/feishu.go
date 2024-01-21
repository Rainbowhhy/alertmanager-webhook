package model

type CardMsg struct {
	MsgType string `json:"msg_type"`
	Card    Card   `json:"card"`
}

type Card struct {
	Elements []Element `json:"elements"`
	Header   Header    `json:"header"`
}

type Element struct {
	Tag  string `json:"tag"`
	Text Body   `json:"text"`
}

type Header struct {
	Template string `json:"template"`
	Title    Body   `json:"title"`
}

type Body struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

func NewCardMsg(color, title string) *CardMsg {
	return &CardMsg{
		MsgType: "interactive",
		Card: Card{
			Header: Header{
				Template: color,
				Title: Body{
					Tag:     "plain_text",
					Content: title,
				},
			},
		},
	}
}

func (c *CardMsg) AddElement(content string) {
	element := Element{
		Tag: "div",
		Text: Body{
			Content: content,
			Tag:     "lark_md",
		},
	}
	c.Card.Elements = append(c.Card.Elements, element)
}
