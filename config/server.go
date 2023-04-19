package config

import "github.com/itxiaolin/metisAi-wechat/internal/core/logger"

type Config struct {
	System  System      `mapstructure:"system" json:"system" yaml:"system"`
	Logger  logger.Conf `mapstructure:"logger" yaml:"logger"`
	OpenAI  OpenAI      `mapstructure:"open-ai" yaml:"open-ai"`
	WxRobot WxRobot     `mapstructure:"wx-robot" yaml:"wx-robot"`
}

type System struct {
	AppName string `mapstructure:"appName" json:"appName" yaml:"appName"` // 应用名称
	PidFile string `mapstructure:"pidFile" json:"pidFile" yaml:"pidFile"` // pid存放位置
}

type OpenAI struct {
	BaseURL string `mapstructure:"base-url" yaml:"base-url"`
	ApiKey  string `mapstructure:"api-key" yaml:"api-key"`
}

type WxRobot struct {
	AutoPass        bool          `mapstructure:"auto-pass" yaml:"auto-pass"`                 //是否自动通过好友请求
	StoragePath     string        `mapstructure:"storage-path" yaml:" storage-path"`          //热登录缓存位置
	SessionTimeout  int64         `mapstructure:"session-timeout" yaml:"session-timeout"`     //上下文过期时间
	ContextCacheNum int           `mapstructure:"context-cache-num" yaml:"context-cache-num"` //上下文缓存数量
	RetryNum        int           `mapstructure:"retry-num" yaml:"retry-num"`                 //chatgpt请求异常重试次数
	ChatGPTModel    string        `mapstructure:"chatGPT-model" yaml:"chatGPT-model"`         //ChatGPT模型
	KeywordPrefix   KeywordPrefix `mapstructure:"keyword-prefix" yaml:"keyword-prefix"`
	Voice           Voice         `mapstructure:"voice" yaml:"voice"`
}

type KeywordPrefix struct {
	ImagePrompt  string `mapstructure:"image-prompt" yaml:"image-prompt"`   //生成图片
	SystemRole   string `mapstructure:"system-role" yaml:"system-role"`     //设置系统角色
	SystemHelp   string `mapstructure:"system-help" yaml:"system-help"`     //指令帮助
	ResetContext string `mapstructure:"reset-context" yaml:"reset-context"` //重置上下文
}

type Voice struct {
	VoiceDir string `mapstructure:"voice-dir" yaml:"voice-dir"` //录音文件临时目录
}
