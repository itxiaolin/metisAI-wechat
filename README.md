# 简介
基于ChatGPT的微信聊天机器人，通过 ChatGPT 接口生成对话内容，使用 openwechat 实现微信消息的接收和自动回复。

- 支持上下文语境的对话。
- 支持重置上下文语境，通过关键词(reset)重置对话上下文语境。
- 支持在群聊@你的机器人 🤖，@机器人即可收到回复。

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
  api-key: "你的api_key"
  base-url: https://api.openai.com/v1
# 支持更换自己指定的url

wx-robot:
  auto-pass: true
  session-timeout: 900
  storage-path: config/storage.json
  retry-num: 3
```

## 命令使用教程
- ./openai-wechat(文件名) start  启动应用
- ./openai-wechat(文件名) start -d 后台启动应用
- ./openai-wechat(文件名) --config config/config.yaml start -d 指定配置文件启动
- ./openai-wechat(文件名) status 查看应用启动情况
- ./openai-wechat(文件名) stop 关闭应用
- ./openai-wechat(文件名) restart 重启应用