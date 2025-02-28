package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"sort"
	"wechat-ai/internal/config"
)

func WechatCheck(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	signature := query.Get("signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")
	echostr := query.Get("echostr")

	sl := []string{config.Wechat.Token, timestamp, nonce}
	sort.Strings(sl)
	sum := sha1.Sum([]byte(sl[0] + sl[1] + sl[2]))
	ok := signature == hex.EncodeToString(sum[:])

	if ok {
		w.Write([]byte(echostr))
		return
	}

	log.Println("此接口为公众号验证，不应该被手动调用，公众号接入校验失败")
}
