package wechat

import (
	"encoding/xml"
	"net/http"
	"time"
)

type Msg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	Content      string   `xml:"Content"`
	Recognition  string   `xml:"Recognition"`

	MsgId int64 `xml:"MsgId,omitempty"`
}

func ParseMsg(data []byte) *Msg {
	var msg Msg
	if err := xml.Unmarshal(data, &msg); err != nil {
		return nil
	}
	return &msg
}

func (msg *Msg) GenerateEchoData(s string) []byte {
	data := Msg{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      s,
	}
	bs, _ := xml.Marshal(&data)
	return bs
}

var success = []byte("success")

func EchoSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(success)
}

func (msg *Msg) EchoText(w http.ResponseWriter, s string) {
	data := Msg{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      s,
	}
	bs, _ := xml.Marshal(&data)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
