package user

import (
	"github.com/itxiaolin/metisAi-wechat/internal/constants"
	"github.com/itxiaolin/metisAi-wechat/internal/domains"
	"github.com/itxiaolin/metisAi-wechat/internal/global"
	"github.com/patrickmn/go-cache"
	"github.com/sashabaranov/go-openai"
	"sync"
	"time"
)

// UserService 用户业务接口
type UserService interface {
	GetUserSessionContext(userId string) *domains.ChatCompletionMessageQueue
	SetUserSessionContext(userId string, question, reply string)
	ClearUserSessionContext(userId string)
	BuildMessages(userId, systemContent, question string) (messages []openai.ChatCompletionMessage)
	SetSystemRole(key, systemRole string)
	GetSystemRole(key string) string
}

var userOnce sync.Once
var userInstance UserService

func Instance() UserService {
	userOnce.Do(func() {
		userInstance = &UserServiceImpl{
			cache: cache.New(time.Second*time.Duration(global.Config.WxRobot.SessionTimeout), time.Minute*10)}
	})
	return userInstance
}

type UserServiceImpl struct {
	cache *cache.Cache
}

// ClearUserSessionContext 清空GTP上下文，接收文本中包含`我要问下一个问题`，并且Unicode 字符数量不超过20就清空
func (s *UserServiceImpl) ClearUserSessionContext(userId string) {
	s.cache.Delete(constants.KEY_USER_SESSION + userId)
}

// GetUserSessionContext 获取用户会话上下文文本
func (s *UserServiceImpl) GetUserSessionContext(userId string) *domains.ChatCompletionMessageQueue {
	sessionContext, ok := s.cache.Get(constants.KEY_USER_SESSION + userId)

	if !ok {
		if global.Config.WxRobot.ContextCacheNum > 0 {
			return domains.NewFixedQueue(global.Config.WxRobot.ContextCacheNum)
		}
		return domains.NewFixedQueue(10)
	}
	queue := sessionContext.(*domains.ChatCompletionMessageQueue)
	return queue
}

// SetUserSessionContext 设置用户会话上下文文本，question用户提问内容，GTP回复内容
func (s *UserServiceImpl) SetUserSessionContext(userId string, question, reply string) {
	messageQueue := s.GetUserSessionContext(constants.KEY_USER_SESSION + userId)
	userMessages := openai.ChatCompletionMessage{
		Role:    constants.CHATGPTRoleUser,
		Content: question,
	}
	messageQueue.Push(userMessages)
	assistantMessages := openai.ChatCompletionMessage{
		Role:    constants.CHATGPTRoleAssistant,
		Content: reply,
	}
	messageQueue.Push(assistantMessages)
	s.cache.Set(userId, messageQueue, time.Second*time.Duration(global.Config.WxRobot.SessionTimeout))
}

func (s *UserServiceImpl) BuildMessages(userId, systemContent, question string) (messages []openai.ChatCompletionMessage) {
	messageQueue := s.GetUserSessionContext(constants.KEY_USER_SESSION + userId)
	systemRoleMessages := openai.ChatCompletionMessage{
		Role:    constants.CHATGPTRoleSystem,
		Content: systemContent,
	}
	messages = append(messages, systemRoleMessages)
	for i := 0; i < messageQueue.Len(); i++ {
		message := messageQueue.Get(i)
		messages = append(messages, message)
	}
	userMessages := openai.ChatCompletionMessage{
		Role:    constants.CHATGPTRoleUser,
		Content: question,
	}
	messages = append(messages, userMessages)
	return
}

func (s *UserServiceImpl) SetSystemRole(userId, systemRole string) {
	s.cache.Set(constants.KEY_SYSTEM_ROLE+userId, systemRole, time.Second*time.Duration(global.Config.WxRobot.SessionTimeout))
	s.cache.Delete(constants.KEY_USER_SESSION + userId)
}

func (s *UserServiceImpl) GetSystemRole(userId string) string {
	sessionContext, ok := s.cache.Get(constants.KEY_SYSTEM_ROLE + userId)
	if !ok {
		return ""
	}
	return sessionContext.(string)
}
