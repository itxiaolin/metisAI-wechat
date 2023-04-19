package conversation

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/itxiaolin/metisAi-wechat/internal/core/logger"
	"github.com/itxiaolin/metisAi-wechat/internal/global"
	"os"
	"sync"
	"time"

	"github.com/eatmoreapple/openwechat"
	tpOpenai "github.com/itxiaolin/metisAi-wechat/third_party/openai"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type ConversationService interface {
	WxChatCompletion(index int, msg *openwechat.Message, atText string, messages []openai.ChatCompletionMessage) (string, error)
	ImagesCompletion(msg *openwechat.Message, atText, content string) (err error)
	Voice2TextCompletion(msg *openwechat.Message) (content string, err error)
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
		logger.Error(nil, "WxChatCompletion request err", zap.Error(err), zap.Any("requestBody", requestBody))
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
	if _, err = msg.ReplyText(fmt.Sprintf("%s %s: %s", atText, content, url)); err != nil {
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

func (s *conversationService) Voice2TextCompletion(msg *openwechat.Message) (content string, err error) {
	voiceResp, err := msg.GetVoice()
	if err != nil {
		logger.Error(nil, "获取语音消息异常", zap.Error(err))
	}
	logger.Info(nil, "获取语言消息成功", zap.String("status", voiceResp.Status))
	defer func() { _ = voiceResp.Body.Close() }()
	_ = os.MkdirAll("config/voice", os.FileMode(0755))
	filePath := fmt.Sprintf("config/voice/%s.mp3", msg.MsgId)
	err = msg.SaveFileToLocal(filePath)
	if err != nil {
		logger.Error(nil, "将语言保存到文件异常", zap.Error(err))
		return "", err
	}
	c := tpOpenai.Instance()
	ctx := context.Background()
	req := openai.AudioRequest{
		Model:       openai.Whisper1,
		FilePath:    filePath,
		Temperature: 0.5,
	}
	TranscriptionResp, err := c.CreateTranscription(ctx, req)
	_ = os.Remove(filePath)
	if err != nil {
		logger.Error(ctx, "通过openai的whisper-1模型转化语言失败", zap.Error(err))
		return
	}
	content = TranscriptionResp.Text
	logger.Info(ctx, "通过openai的whisper-1模型转化语言成功", zap.String("内容", content))
	return
}
