package openai

import (
	"github.com/itxiaolin/metisAi-wechat/internal/global"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

var once sync.Once
var instance *openai.Client

func Instance() *openai.Client {
	once.Do(func() {
		conf := openai.DefaultConfig(global.Config.OpenAI.ApiKey)
		conf.BaseURL = global.Config.OpenAI.BaseURL
		instance = openai.NewClientWithConfig(conf)
	})
	return instance
}
