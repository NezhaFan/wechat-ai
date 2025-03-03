package handler

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"wechat-ai/internal/config"
	"wechat-ai/internal/service/model"
	"wechat-ai/internal/service/wechat"
)

// https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Passive_user_reply_message.html
// 微信服务器在五秒内收不到响应会断掉连接，并且重新发起请求，总共重试三次
func ReceiveMsg(w http.ResponseWriter, r *http.Request) {
	bs, _ := io.ReadAll(r.Body)
	msg := wechat.ParseMsg(bs)

	if msg == nil {
		log.Println("xml格式公众号消息接口，请勿手动调用")
		wechat.EchoSuccess(w)
		return
	}

	// 非文本不回复(返回success表示不回复)
	switch msg.MsgType {
	// 未写的类型
	default:
		log.Printf("未实现的消息类型%s\n", msg.MsgType)
		wechat.EchoSuccess(w)
	case "event":
		switch msg.Event {
		default:
			log.Printf("未实现的事件%s\n", msg.Event)
			wechat.EchoSuccess(w)
		case "subscribe":
			msg.EchoText(w, config.Wechat.SubscribeMsg)
			return
		case "unsubscribe":
			log.Println("取消关注:", msg.FromUserName)
			wechat.EchoSuccess(w)
			return
		}
	// https://developers.weixin.qq.com/community/minihome/doc/0004826962c5c81c0540cb9e365401?page=1
	case "voice":
		msg.EchoText(w, "不好意思哈，我听不到语音消息～")
	case "text":

	}

	ch := GetUserChan(msg)

	select {
	// 前两次超时不回答
	case <-time.After(time.Second * 5):
		// log.Println("5s超时")
	case result := <-ch:
		msg.EchoText(w, result)
	}

}

var (
	replyCache sync.Map
)

// 使用chan的目的在于能提前返回
func GetUserChan(msg *wechat.Msg) (ch chan string) {
	replyCh, ok := replyCache.Load(msg.MsgId)
	if !ok {
		ch = make(chan string, 1)
		replyCache.Store(msg.MsgId, ch)

		go func(msgid int64) {
			resultCh := make(chan string)
			go func() {
				resultCh <- model.Chat(msg.FromUserName, msg.Content)
			}()

			select {
			case <-time.After(time.Second * 14):
				ch <- "抱歉，无法在微信限制时间内做出应答"
			case reply := <-resultCh:
				ch <- reply
			}
			close(ch)
			replyCache.Delete(msgid)
		}(msg.MsgId)
	} else {
		ch = replyCh.(chan string)
	}
	return
}
