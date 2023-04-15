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
	Instance()
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "记住你是图片扩散模型dell e2的prompt专家，每当我给你提需求的时候，用英文回复你梳理好的prompt语句，给一条最好最精简的",
			}, {
				Role:    "user",
				Content: "生成两个美女在太空船喝咖啡的图片",
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
