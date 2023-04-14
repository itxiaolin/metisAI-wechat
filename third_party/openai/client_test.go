package openai

import (
	"context"
	"fmt"
	"github.com/itxiaolin/openai-wechat/internal/global"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func TestImagesCompletion(t *testing.T) {
	global.Config.OpenAI.ApiKey = ""
	global.Config.OpenAI.BaseURL = ""
	requestBody := openai.ImageRequest{
		Prompt: "生成一只白色的猫",
		N:      1,
		Size:   "1024x1024",
	}
	response, err := Instance().CreateImage(context.Background(), requestBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	url := response.Data[0].URL
	fmt.Println(url)
}

func TestChatCompletion(t *testing.T) {
	global.Config.OpenAI.ApiKey = ""
	global.Config.OpenAI.BaseURL = ""
	instance = Instance()
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{{
			Role:    "user",
			Content: "你好",
		},
		},
	}
	response, err := instance.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
}
