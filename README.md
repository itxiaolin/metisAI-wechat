# 简介
基于ChatGPT的微信聊天机器人，通过 ChatGPT 接口生成对话内容，使用 openwechat 实现微信消息的接收和自动回复。

# 快速开始
支持 Linux、MacOS、Windows 系统（可在Linux服务器上长期运行),不需安装安装任何环境,如果是本地运行，需要安装golang

## 命令使用教程
- openai-wechat start  启动应用
- openai-wechat start -d 后台启动引用
- openai-wechat --config config.yaml start -d 指定配置文件启动
- openai-wechat status 查看应用启动情况
- openai-wechat stop 关闭应用
- openai-wechat restart 重启应用