package user

import (
	"github.com/itxiaolin/openai-wechat/internal/constants"
	"github.com/itxiaolin/openai-wechat/internal/domains"
	"github.com/itxiaolin/openai-wechat/internal/global"
	"github.com/patrickmn/go-cache"
	"github.com/sashabaranov/go-openai"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// UserService 用户业务接口
type UserService interface {
	GetUserSessionContext(userId string) *domains.ChatCompletionMessageQueue
	SetUserSessionContext(userId string, question, reply string)
	ClearUserSessionContext(userId string, msg string) bool
	BuildMessages(userId, systemContent, question string) (messages []openai.ChatCompletionMessage)
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
func (s *UserServiceImpl) ClearUserSessionContext(userId string, msg string) bool {
	if strings.Contains(msg, "清除上下文") && utf8.RuneCountInString(msg) < 20 {
		s.cache.Delete(userId)
		return true
	}
	return false
}

// GetUserSessionContext 获取用户会话上下文文本
func (s *UserServiceImpl) GetUserSessionContext(userId string) *domains.ChatCompletionMessageQueue {
	sessionContext, ok := s.cache.Get(userId)
	if !ok {
		return domains.NewFixedQueue(10)
	}
	queue := sessionContext.(*domains.ChatCompletionMessageQueue)
	return queue
}

// SetUserSessionContext 设置用户会话上下文文本，question用户提问内容，GTP回复内容
func (s *UserServiceImpl) SetUserSessionContext(userId string, question, reply string) {
	messageQueue := s.GetUserSessionContext(userId)
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
	messageQueue := s.GetUserSessionContext(userId)
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
