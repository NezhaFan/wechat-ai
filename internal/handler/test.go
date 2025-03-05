package handler

import (
	"net/http"
	"wechat-ai/internal/service/model"
)

func Test(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")

	reply := model.Chat("test", msg)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(reply))
}
