package config

import "github.com/itxiaolin/openai-wechat/internal/core/logger"

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
	AutoPass       bool   `mapstructure:"autoPass" yaml:"autoPass"`
	SessionTimeout int64  `mapstructure:"sessionTimeout" yaml:"sessionTimeout"`
	StoragePath    string `mapstructure:"storagePath" yaml:"storagePath"`
	RetryNum       int    `mapstructure:"retryNum" yaml:"retryNum"`
}
