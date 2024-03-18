package model

import "encoding/xml"

const MsgTypeEvent = "event"
const MsgTypeText = "text"
const MsgTypeNews = "news"

const EventTypeSubscribe = "subscribe"
const EventTypeUnsubscribe = "unsubscribe"

// ReqMessage 接收消息结构体
type ReqMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Event        string
	Content      string
	MsgId        int64
}

// RespMessage 响应消息基础结构体
type RespMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}

// RespTextMessage 文本响应消息结构体
type RespTextMessage struct {
	RespMessage
	Content string
}

// RespNewsMessage 图文响应消息结构体
type RespNewsMessage struct {
	RespMessage
	ArticleCount int
	Articles     []ArticleItem
}

// ArticleItem 图文结构体
type ArticleItem struct {
	Article Article
}

// Article 图文结构体
type Article struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
	XMLName     xml.Name `xml:"item"`
}
