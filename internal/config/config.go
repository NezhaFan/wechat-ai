package config

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	Port string

	LLM struct {
		Api   string
		Key   string
		Model string

		Prompt      string
		Temperature float32
		MaxTokens   uint32
		History     int16
	}

	Wechat struct {
		Token        string
		SubscribeMsg string
	}
)

func init() {
	// 使用flag接收 config文件
	configFile := flag.String("c", "./config.yaml", "配置文件名")
	flag.Parse()

	// 读取配置
	viper.SetConfigFile(*configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("解析配置文件config.yaml失败:", err.Error())
		os.Exit(0)
	}

	viper.UnmarshalKey("port", &Port)
	viper.UnmarshalKey("llm", &LLM)
	viper.UnmarshalKey("wechat", &Wechat)

	if LLM.Key == "" || LLM.Api == "" || LLM.Model == "" || LLM.MaxTokens <= 0 {
		log.Println("大模型配置错误")
		os.Exit(0)
	}

	if Wechat.Token == "" {
		log.Println("未设置公众号token，公众号功能不可用")
	}

}
