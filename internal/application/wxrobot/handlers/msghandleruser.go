package handlers

import (
	"fmt"
	"github.com/itxiaolin/openai-wechat/internal/application/wxrobot/services/conversation"
	"github.com/itxiaolin/openai-wechat/internal/application/wxrobot/services/user"
	"github.com/itxiaolin/openai-wechat/internal/core/logger"
	"github.com/itxiaolin/openai-wechat/internal/global"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	if msg.IsVoice() {
		return g.ReplyVoice(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	sender, err := msg.Sender()
	self, _ := msg.Bot().GetCurrentUser()
	logger.Debug(nil, fmt.Sprintf("Received User %v Text Msg : %v", sender.NickName, msg.Content))
	content := strings.TrimSpace(msg.Content)
	contextKey := self.ID() + "-" + sender.UserName
	requestText := strings.Trim(strings.TrimSpace(msg.Content), "\n")
	if requestText == "" {
		_, err = msg.ReplyText("很抱歉，我不太明白您想表达什么，请您再提供更多的信息或者问题，让我能够更好地为您服务。")
		return nil
	}
	if user.Instance().ClearUserSessionContext(contextKey, content) {
		_, err = msg.ReplyText("上下文已经清空了，你可以问下一个问题啦。")
		if err != nil {
			logger.Error(nil, "response user error", zap.Error(err))
		}
		return nil
	}
	if global.Config.WxRobot.RobotKeywordPrompt.ImagePrompt != "" && strings.HasPrefix(requestText, global.Config.WxRobot.RobotKeywordPrompt.ImagePrompt) {
		return conversation.Instance().ImagesCompletion(msg, "", content)
	}
	err = g.replyMsg(msg, contextKey, content)
	if err != nil {
		return err
	}
	return nil
}

func (g *UserMessageHandler) ReplyVoice(msg *openwechat.Message) error {
	content, err := conversation.Instance().Voice2TextCompletion(msg)
	if err != nil {
		return err
	}
	sender, _ := msg.Sender()
	self, _ := msg.Bot().GetCurrentUser()
	contextKey := self.ID() + "-" + sender.UserName
	err = g.replyMsg(msg, contextKey, content)
	if err != nil {
		return err
	}
	return nil
}

func (g *UserMessageHandler) replyMsg(msg *openwechat.Message, contextKey string, content string) error {
	completionMessages := user.Instance().BuildMessages(contextKey, systemContent, content)
	reply, err := conversation.Instance().WxChatCompletion(global.Config.WxRobot.RetryNum, msg, "", completionMessages)
	if err != nil || reply == "" {
		logger.Error(nil, "gpt request error", zap.Error(err))
		_, err = msg.ReplyText("机器人挂了，我一会儿修复一下。")
		if err != nil {
			logger.Error(nil, "response user error", zap.Error(err))
		}
		return err
	}
	user.Instance().SetUserSessionContext(contextKey, content, reply)
	if err != nil {
		logger.Error(nil, "response user error", zap.Error(err))
		return err
	}
	return nil
}
