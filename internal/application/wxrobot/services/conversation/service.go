package conversation

import (
	"bytes"
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/itxiaolin/openai-wechat/internal/core/logger"
	"github.com/itxiaolin/openai-wechat/internal/global"
	"sync"
	"time"

	"github.com/eatmoreapple/openwechat"
	tpOpenai "github.com/itxiaolin/openai-wechat/third_party/openai"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type ConversationService interface {
	WxChatCompletion(index int, msg *openwechat.Message, atText string, messages []openai.ChatCompletionMessage) (string, error)
	ImagesCompletion(msg *openwechat.Message, atText, content string) (err error)
}

var wxChatOnce sync.Once
var wxChatInstance ConversationService

func Instance() ConversationService {
	wxChatOnce.Do(func() {
		wxChatInstance = &conversationService{}
	})
	return wxChatInstance
}

type conversationService struct {
}

func (s *conversationService) WxChatCompletion(index int, msg *openwechat.Message, atText string, messages []openai.ChatCompletionMessage) (string, error) {
	ctx := context.Background()
	var model string
	if global.Config.WxRobot.ChatGPTModel != "" {
		model = global.Config.WxRobot.ChatGPTModel
	} else {
		model = openai.GPT3Dot5Turbo
	}
	requestBody := openai.ChatCompletionRequest{
		Model: model,
	}
	requestBody.Messages = messages
	response, err := tpOpenai.Instance().CreateChatCompletion(ctx, requestBody)
	if err != nil {
		if index > 0 {
			time.Sleep(500 * time.Microsecond)
			return s.WxChatCompletion(index-1, msg, atText, messages)
		}
		return "", err
	}
	text := response.Choices[0].Message.Content
	if text != "" {
		replyText := atText + text
		_, _ = msg.ReplyText(replyText)
		logger.Info(nil, "message returned successfully：", zap.String("text", replyText))
	}
	return text, nil
}

func (s *conversationService) ImagesCompletion(msg *openwechat.Message, atText, content string) (err error) {
	c := tpOpenai.Instance()
	ctx := context.Background()
	requestBody := openai.ImageRequest{
		Prompt: content,
		N:      1,
		Size:   "1024x1024",
	}
	response, err := c.CreateImage(ctx, requestBody)
	if err != nil {
		_, err = msg.ReplyText("ChatGPT现在满负荷运转,请稍后重试。")
		logger.Error(ctx, "request chat gpt err", zap.Error(err))
		return err
	}
	url := response.Data[0].URL
	if _, err = msg.ReplyText(atText + "chatGPT生成图片: " + url); err != nil {
		logger.Error(ctx, "response group error", zap.Error(err))
		return err
	}
	resp, err := resty.New().R().Get(url)
	if err == nil && resp.IsSuccess() {
		reader := bytes.NewReader(resp.Body())
		_, _ = msg.ReplyImage(reader)
	}
	return nil
}
