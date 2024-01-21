package model

type Message struct {
	QywechatMessage QyWechatMessage
	FeishuMessage   FeiShuMessage
	DingdingMessage DingDingMessage
}

type QyWechatMessage struct {
	MarkdownFiring   *QyWeChatMarkdown
	MarkdownResolved *QyWeChatMarkdown
}

type FeiShuMessage struct {
	CardFiring   *CardMsg
	CardResolved *CardMsg
}

type DingDingMessage struct {
	DmarkdownFiring   *DingDingMarkdown
	DmarkdownResolved *DingDingMarkdown
}

func NewMessage(markdownFiring, markdownResolved *QyWeChatMarkdown, cardFiring, cardResolved *CardMsg, dmarkdownFiring, dmarkdownResolved *DingDingMarkdown) *Message {
	return &Message{
		QywechatMessage: QyWechatMessage{
			MarkdownFiring:   markdownFiring,
			MarkdownResolved: markdownResolved,
		},
		FeishuMessage: FeiShuMessage{
			CardFiring:   cardFiring,
			CardResolved: cardResolved,
		},
		DingdingMessage: DingDingMessage{
			DmarkdownFiring:   dmarkdownFiring,
			DmarkdownResolved: dmarkdownResolved,
		},
	}
}
