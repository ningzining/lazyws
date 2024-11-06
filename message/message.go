package message

// TopicType 主题类型
type TopicType string

const (
	SubscribeTopicType   TopicType = "subscribe"   // 消息订阅
	UnsubscribeTopicType TopicType = "unsubscribe" // 取消订阅
)

type Request struct {
	TopicType TopicType `json:"topic_type"` // 消息类型
	Topic     string    `json:"topic"`      // 消息主题
}

// MessageType 消息类型
type MessageType string

const (
	NoticeType MessageType = "notice" // 通知
	WarnType   MessageType = "warn"   // 警告
	ErrorType  MessageType = "error"  // 错误

	InfoType MessageType = "info" // 信息
)

type Response struct {
	MessageType MessageType `json:"message_type"` // 消息类型
	Payload     any         `json:"payload"`      // 消息内容，根据消息类型解析对应的结构
}
