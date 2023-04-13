package conversation

import (
	"context"
	"github.com/itxiaolin/openai-wechat/internal/core/logger"
	"sync"
	"time"

	"github.com/eatmoreapple/openwechat"
	tpOpenai "github.com/itxiaolin/openai-wechat/third_party/openai"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type ConversationService interface {
	WxChatCompletion(index int, msg *openwechat.Message, atText string, messages []openai.ChatCompletionMessage) (string, error)
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
	c := tpOpenai.Instance()
	ctx := context.Background()
	requestBody := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
	}
	requestBody.Messages = messages
	response, err := c.CreateChatCompletion(ctx, requestBody)
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
