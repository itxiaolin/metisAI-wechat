package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/itxiaolin/openai-wechat/internal/core/logger"
	"github.com/itxiaolin/openai-wechat/internal/global"
	"go.uber.org/zap"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

var (
	BotLoginTimeMap = make(map[string]int64)
)

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	loginTime := BotLoginTimeMap[msg.Bot().UUID()]
	if msg.CreateTime < loginTime {
		logger.Debug(nil, "历史消息不处理", zap.String("msg", msg.Content))
		return
	}
	logger.Info(nil, "handle Received", zap.String("msg", msg.Content))

	if msg.IsSendByGroup() {
		go handlers[GroupHandler].handle(msg)
		return
	}

	if msg.IsFriendAdd() {
		if global.Config.WxRobot.AutoPass {
			_, err := msg.Agree(GetHelpText(""))
			if err != nil {
				logger.Error(nil, "add friend agree", zap.Error(err))
				return
			}
		}
	}

	go handlers[UserHandler].handle(msg)
}
