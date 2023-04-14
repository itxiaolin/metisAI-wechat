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

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
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
	sender, err := msg.Sender()
	group := openwechat.Group{User: sender}
	self, _ := msg.Bot().GetCurrentUser()
	logger.Debug(nil, "Received msg", zap.String("Group", group.NickName), zap.String("Text Msg", msg.Content))
	if !msg.IsAt() {
		return nil
	}

	groupSender, err := msg.SenderInGroup()
	if err != nil {
		logger.Error(nil, "get sender in group", zap.Error(err))
		return err
	}
	atText := "@" + groupSender.NickName + " "
	contextKey := self.ID() + "-" + group.NickName + "-" + groupSender.UserName
	content := strings.TrimSpace(msg.Content)
	if user.Instance().ClearUserSessionContext(contextKey, content) {
		_, err = msg.ReplyText(atText + "上下文已经清空了，你可以问下一个问题啦。")
		if err != nil {
			logger.Error(nil, "response user error", zap.Error(err))
		}
		return nil
	}
	msgContent := strings.TrimSpace(strings.ReplaceAll(content, fmt.Sprintf("@%s", self.NickName), ""))
	if msgContent == "" {
		return nil
	}
	if global.Config.WxRobot.RobotKeywordPrompt.ImagePrompt != "" && strings.HasPrefix(msgContent, global.Config.WxRobot.RobotKeywordPrompt.ImagePrompt) {
		return conversation.Instance().ImagesCompletion(msg, atText, content)
	}
	completionMessages := user.Instance().BuildMessages(contextKey, systemContent, msgContent)
	logger.Info(nil, fmt.Sprintf("request chatGPT by group: %v, NickName:%s ,requestText: %v",
		group.NickName, groupSender.NickName, msgContent), zap.String("UserName", groupSender.UserName))
	reply, err := conversation.Instance().WxChatCompletion(global.Config.WxRobot.RetryNum, msg, atText, completionMessages)
	if err != nil {
		logger.Error(nil, "gpt request error", zap.Error(err))
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
	user.Instance().SetUserSessionContext(contextKey, msgContent, reply)
	return err
}
