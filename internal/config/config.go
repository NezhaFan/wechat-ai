package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	Http struct {
		Port  string `json:"port"`
		Proxy string `json:"proxy"`
	}

	OpenAI struct {
		Key    string `json:"key"`
		Params struct {
			Api    string `json:"api"`
			Model  string `json:"model"`
			Prompt string `json:"prompt"`
		} `json:"params"`
	}

	Wechat struct {
		Token string `json:"token"`
	}
)

func init() {

	// 读取配置
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("解析配置文件config.yaml失败:", err.Error())
		os.Exit(0)
	}

	viper.UnmarshalKey("http", &Http)
	viper.UnmarshalKey("openai", &OpenAI)
	viper.UnmarshalKey("wechat", &Wechat)

	if OpenAI.Key == "" {
		fmt.Println("OpenAI的Key不能为空")
		os.Exit(0)
	}

	if Wechat.Token == "" {
		fmt.Println("未设置公众号token，公众号功能不可用")
	}

}
