package handlers

import (
	"fmt"
	"github.com/itxiaolin/metisAi-wechat/internal/application/wxrobot/services/conversation"
	"github.com/itxiaolin/metisAi-wechat/internal/application/wxrobot/services/user"
	"github.com/itxiaolin/metisAi-wechat/internal/core/logger"
	"github.com/itxiaolin/metisAi-wechat/internal/global"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsTickledMe() {
		_, _ = msg.ReplyText("我在，有什麼我可以幫你的地方嗎？")
		return nil
	}
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	sender, _ := msg.Sender()
	group := openwechat.Group{User: sender}
	self, _ := msg.Bot().GetCurrentUser()
	logger.Debug(nil, fmt.Sprintf("Received group"), zap.String("NickName", group.NickName),
		zap.String("UserName", group.UserName), zap.String("id", group.ID()))
	if !msg.IsAt() {
		return nil
	}
	groupSender, _ := msg.SenderInGroup()
	logger.Debug(nil, fmt.Sprintf("Received groupSender"), zap.String("NickName", groupSender.NickName),
		zap.String("UserName", groupSender.UserName), zap.String("id", groupSender.ID()))
	atText := "@" + groupSender.NickName + " "
	contextKey := self.ID() + "-" + group.ID() + "-" + groupSender.UserName
	requestText := strings.TrimSpace(strings.ReplaceAll(strings.TrimSpace(msg.Content), fmt.Sprintf("@%s", self.NickName), ""))
	if requestText == "" {
		_, _ = msg.ReplyText(GetHelpText(atText))
		return nil
	}
	if isKeywordPrefix, err := KeywordPrefixHandlerInstance().Handle(msg, atText, requestText, contextKey); isKeywordPrefix {
		return err
	}
	completionMessages := user.Instance().BuildMessages(contextKey, user.Instance().GetSystemRole(contextKey), requestText)
	logger.Info(nil, fmt.Sprintf("request chatGPT by group: %v, NickName:%s ,requestText: %v",
		group.NickName, groupSender.NickName, requestText), zap.String("UserName", groupSender.UserName))
	reply, err := conversation.Instance().WxChatCompletion(global.Config.WxRobot.RetryNum, msg, atText, completionMessages)
	if err != nil {
		logger.Error(nil, "chat completion request error", zap.Error(err))
		_, err = msg.ReplyText("机器人神了，我一会发现了就去修。")
		if err != nil {
			logger.Error(nil, "response group error", zap.Error(err))
		}
		return err
	}
	if reply == "" {
		_, err = msg.ReplyText(atText + "ChatGPT现在满负荷运转,请稍后重试。")
		return nil
	}
	user.Instance().SetUserSessionContext(contextKey, requestText, reply)
	return err
}
