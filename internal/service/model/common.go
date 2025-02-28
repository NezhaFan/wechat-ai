package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"wechat-ai/internal/config"
)

type Request struct {
	Model       string           `json:"model"`
	Messages    []RequestMessage `json:"messages"`
	Temperature float32          `json:"temperature"`
	MaxTokens   uint32           `json:"max_tokens"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	ID      string `json:"id"`
	Choices []struct {
		// 流式
		// Delta struct {
		// 	Content string `json:"content"`
		// } `json:"delta,omitempty"`
		// json
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message,omitempty"`
		// json
		Usage struct {
			PromptTokens         uint16 `json:"prompt_tokens"`
			PromptCacheHitTokens uint16 `json:"prompt_cache_hit_tokens"`
			CompletionTokens     uint16 `json:"completion_tokens"`
		} `json:"usage,omitempty"`
	} `json:"choices"`
}

func Chat(uid string, msg string) string {
	messages := getMessages(uid)

	// 第一次说话，加上预设的prompt
	if len(messages) == 0 && config.LLM.Prompt != "" {
		addMessages(uid, "system", config.LLM.Prompt)
		messages = append(messages, RequestMessage{Role: "system", Content: config.LLM.Prompt})
	}

	if len(messages) > 0 {
		if config.LLM.History <= 0 {
			messages = messages[:1]
		} else {
			n := int(config.LLM.History * 2)
			size := len(messages)
			if n < size-1 {
				copy(messages[1:], messages[size-n:])
				fmt.Println("变为", messages)
				messages = messages[:n+1]
			}
		}
	}

	messages = append(messages, RequestMessage{Role: "user", Content: msg})

	fmt.Println(messages)
	req := Request{
		Model:       config.LLM.Model,
		Messages:    messages,
		Temperature: config.LLM.Temperature,
		MaxTokens:   config.LLM.MaxTokens,
	}
	data, _ := json.Marshal(&req)

	url := strings.TrimSuffix(config.LLM.Api, "/") + "/chat/completions"
	echo, err := post(url, config.LLM.Key, data)
	if err != nil {
		log.Println("发起请求失败：" + err.Error())
		log.Println("请求失败数据：" + string(data))
		return "抱歉，网络错误"
	}
	var res Response
	err = json.Unmarshal(echo, &res)
	if err != nil {
		return err.Error()
	}

	reply := string(res.Choices[0].Message.Content)
	addMessages(uid, "user", msg)
	addMessages(uid, "assistant", reply)

	return reply
}

func post(url, auth string, data []byte) (body []byte, err error) {
	client := http.Client{Timeout: time.Second * 50}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+auth)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("请求发起失败，错误码：" + res.Status)
	}
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
