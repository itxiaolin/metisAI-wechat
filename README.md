# 简介
基于ChatGPT的微信聊天机器人，通过 ChatGPT 接口生成对话内容，使用 openwechat 实现微信消息的接收和自动回复。

- 上下文语境: 支持重置上下文语境，通过关键词(清除上下文)重置对话上下文语境。
- 群聊机器人: 支持在群聊@你的机器人 🤖，@机器人即可收到回复。
- 角色扮演: 支持自定义chatGPT的system角色，可实现多种风格的机器人
- 图片生成: 支持根据描述生成图片，默认前缀/image-prompt，支持修改配置
- 语音识别: 支持接收和处理语音消息，通过文字回复

# 快速开始
支持 Linux、MacOS、Windows 系统（可在Linux服务器上长期运行),不需安装安装任何环境,如果是本地代码运行，需要安装golang

```shell
# 获取项目
git clone https://github.com/itxiaolin/openai-wechat.git
# 进入项目目录
cd openai-wechat
# 修改配置(配置api_key)
open-ai:
  api-key: "你的api_key"
  base-url: https://api.openai.com/v1
# 依赖下载
go mod tidy 
# 启动项目
go run main.go
```

## 默认配置
```yaml
system:
  appName: openai-wechat
  pidFile: config/bin/openai-wechat.lock

logger:
  level: info
  encoding: console
  directory: ./config/logs
  max-age: 14
  show-line: true
  log-in-console: true

open-ai:
  api-key: ""
  base-url: ""

wx-robot:
  auto-pass: true
  session-timeout: 900
  storage-path: config/storage.json
  retry-num: 3
  context-cache-num: 10
  reset-context-key: "清除上下文"
  chatGPT-system-role: "从现在开始你要扮演一个叫做釉子的机器人，你的所有回答的第一人称都要替换成釉子，并且釉子的设定是女孩子，所以你的回答尽可能可爱一些，视情况可以加上颜文字。"
  chatGPT-model: "gpt-3.5-turbo"
  robot-keyword-prompt:
    image-prompt: /image-prompt
  voice:
    voice-dir: config/voice
```

## 命令使用教程
- ./main start  启动应用
- ./main start -d 后台启动应用
- ./main --config config/config.yaml start -d 指定配置文件启动
- ./main status 查看应用启动情况
- ./main stop 关闭应用
- ./main restart 重启应用