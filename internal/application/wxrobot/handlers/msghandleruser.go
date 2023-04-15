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
	if msg.IsTickledMe() {
		_, _ = msg.ReplyText("我在，有什麼我可以幫你的地方嗎？")
		return nil
	}
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

// ReplyText 私聊回复
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	sender, _ := msg.Sender()
	self, _ := msg.Bot().GetCurrentUser()
	logger.Debug(nil, fmt.Sprintf("Received NickName: %s UserName: %s Text Msg : %v", sender.NickName, sender.UserName, msg.Content),
		zap.String("id", sender.ID()))
	content := strings.TrimSpace(msg.Content)
	contextKey := self.ID() + "-" + sender.UserName
	requestText := strings.Trim(strings.TrimSpace(msg.Content), "\n")
	if requestText == "" {
		_, _ = msg.ReplyText(GetHelpText(""))
		return nil
	}
	if isKeywordPrefix, err := KeywordPrefixHandlerInstance().Handle(msg, "", requestText, contextKey); isKeywordPrefix {
		return err
	}
	if err := g.replyMsg(msg, contextKey, content); err != nil {
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
	completionMessages := user.Instance().BuildMessages(contextKey, user.Instance().GetSystemRole(contextKey), content)
	reply, err := conversation.Instance().WxChatCompletion(global.Config.WxRobot.RetryNum, msg, "", completionMessages)
	if err != nil || reply == "" {
		logger.Error(nil, "chat completion request error", zap.Error(err))
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
