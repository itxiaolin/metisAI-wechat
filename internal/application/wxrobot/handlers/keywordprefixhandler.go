package handlers

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/itxiaolin/metisAi-wechat/internal/application/wxrobot/services/conversation"
	"github.com/itxiaolin/metisAi-wechat/internal/application/wxrobot/services/user"
	"github.com/itxiaolin/metisAi-wechat/internal/core/logger"
	"github.com/itxiaolin/metisAi-wechat/internal/global"
	"go.uber.org/zap"
	"strings"
	"sync"
)

var KeywordPrefixHandlerOnce sync.Once
var keywordPrefixHandlerInstance KeywordPrefixHandler

func KeywordPrefixHandlerInstance() KeywordPrefixHandler {
	KeywordPrefixHandlerOnce.Do(func() {
		resetContextHandler := &ResetContextHandler{}
		helpContextHandler := &HelpContextHandler{}
		imagePromptHandler := &ImagePromptHandler{}
		systemRoleHandler := &SystemRoleHandler{}
		resetContextHandler.SetNext(helpContextHandler)
		helpContextHandler.SetNext(imagePromptHandler)
		imagePromptHandler.SetNext(systemRoleHandler)
		keywordPrefixHandlerInstance = resetContextHandler
	})
	return keywordPrefixHandlerInstance
}

type KeywordPrefixHandler interface {
	SetNext(handler KeywordPrefixHandler)
	Handle(msg *openwechat.Message, atText, requestText, contextKey string) (bool, error)
}

type ResetContextHandler struct {
	nextHandler KeywordPrefixHandler
}

func (r *ResetContextHandler) SetNext(handler KeywordPrefixHandler) {
	r.nextHandler = handler
}

func (r *ResetContextHandler) Handle(msg *openwechat.Message, atText, requestText, contextKey string) (bool, error) {
	if global.Config.WxRobot.KeywordPrefix.ResetContext != "" && strings.HasPrefix(requestText, global.Config.WxRobot.KeywordPrefix.ResetContext) {
		user.Instance().ClearUserSessionContext(contextKey)
		_, err := msg.ReplyText(atText + "上下文已经清空了，你可以问下一个问题啦。")
		if err != nil {
			logger.Error(nil, "response user error", zap.Error(err))
		}
		return true, err
	}
	if r.nextHandler == nil {
		return false, nil
	}
	return r.nextHandler.Handle(msg, atText, requestText, contextKey)
}

type HelpContextHandler struct {
	nextHandler KeywordPrefixHandler
}

func (r *HelpContextHandler) SetNext(handler KeywordPrefixHandler) {
	r.nextHandler = handler
}

func (r *HelpContextHandler) Handle(msg *openwechat.Message, atText, requestText, contextKey string) (bool, error) {
	if global.Config.WxRobot.KeywordPrefix.SystemHelp != "" && strings.HasPrefix(requestText, global.Config.WxRobot.KeywordPrefix.SystemHelp) {
		requestText = strings.TrimPrefix(requestText, global.Config.WxRobot.KeywordPrefix.SystemHelp)
		_, err := msg.ReplyText(GetHelpText(atText))
		if err != nil {
			logger.Error(nil, "response user error", zap.Error(err))
		}
		return true, err
	}
	if r.nextHandler == nil {
		return false, nil
	}
	return r.nextHandler.Handle(msg, atText, requestText, contextKey)
}

type ImagePromptHandler struct {
	nextHandler KeywordPrefixHandler
}

func (r *ImagePromptHandler) SetNext(handler KeywordPrefixHandler) {
	r.nextHandler = handler
}

func (r *ImagePromptHandler) Handle(msg *openwechat.Message, atText, requestText, contextKey string) (bool, error) {
	if global.Config.WxRobot.KeywordPrefix.ImagePrompt != "" && strings.HasPrefix(requestText, global.Config.WxRobot.KeywordPrefix.ImagePrompt) {
		requestText = strings.TrimPrefix(requestText, global.Config.WxRobot.KeywordPrefix.ImagePrompt)
		err := conversation.Instance().ImagesCompletion(msg, atText, requestText)
		if err != nil {
			logger.Error(nil, "response user error", zap.Error(err))
		}
		return true, err
	}
	if r.nextHandler == nil {
		return false, nil
	}
	return r.nextHandler.Handle(msg, atText, requestText, contextKey)
}

type SystemRoleHandler struct {
	nextHandler KeywordPrefixHandler
}

func (r *SystemRoleHandler) SetNext(handler KeywordPrefixHandler) {
	r.nextHandler = handler
}

func (r *SystemRoleHandler) Handle(msg *openwechat.Message, atText, requestText, contextKey string) (bool, error) {
	if global.Config.WxRobot.KeywordPrefix.SystemRole != "" && strings.HasPrefix(requestText, global.Config.WxRobot.KeywordPrefix.SystemRole) {
		requestText = strings.TrimPrefix(requestText, global.Config.WxRobot.KeywordPrefix.SystemRole)
		user.Instance().SetSystemRole(contextKey, requestText)
		_, err := msg.ReplyText(atText + "系统角色设置成功。")
		if err != nil {
			logger.Error(nil, "response user error", zap.Error(err))
		}
		return true, err
	}
	if r.nextHandler == nil {
		return false, nil
	}
	return r.nextHandler.Handle(msg, atText, requestText, contextKey)
}

func GetHelpText(atText string) string {
	return fmt.Sprintf("%s\n"+
		"我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。\n"+
		"@机器人: 你可以在群聊中@我进行交互，我会及时回复哦 🤖。\n"+
		"私聊语音: 发送语音消息，我也能够接收并回复文字消息\n"+
		"关键词: 输入关键词以触发更多功能\n"+
		"%s : 查看帮助提示。\n"+
		"%s : 重置对话上下文语境。\n"+
		"%s : 支持自定义chatGPT的system角色，可指定聊天时的角色\n"+
		"%s : 支持根据描述生成图片\n",
		atText,
		global.Config.WxRobot.KeywordPrefix.SystemHelp,
		global.Config.WxRobot.KeywordPrefix.ResetContext,
		global.Config.WxRobot.KeywordPrefix.SystemRole,
		global.Config.WxRobot.KeywordPrefix.ImagePrompt)
}
