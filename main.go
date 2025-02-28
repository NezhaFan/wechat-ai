package main

import (
	"log"
	"net/http"
	"wechat-ai/bootstrap"
	"wechat-ai/internal/config"
	"wechat-ai/internal/handler"
)

func init() {

}

func main() {
	r := bootstrap.New()

	// 微信消息处理
	r.POST("/wx", handler.ReceiveMsg)
	// 用于公众号自动验证
	r.GET("/wx", handler.WechatCheck)
	// 用于测试 curl "http://127.0.0.1:$PORT/"
	r.GET("/", handler.Test)

	log.Printf("启动服务，测试：curl 'http://127.0.0.1:%s?msg=你好' ", config.Port)
	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		panic(err)
	}
}
