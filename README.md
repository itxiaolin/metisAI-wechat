# 简介
基于ChatGPT的微信聊天机器人，通过 ChatGPT 接口生成对话内容，使用 openwechat 实现微信消息的接收和自动回复。
- [x] **上下文语境：** 支持重置上下文语境，，通过默认前缀/reset-context重置对话上下文语境。
- [x] **群聊机器人：** 支持在群聊@你的机器人 🤖，@机器人即可收到回复。
- [x] **角色扮演：** 支持自定义chatGPT的system角色，通过默认前缀/system-role可指定聊天时的角色
- [x] **图片生成：** 支持根据描述生成图片，默认前缀/image-prompt，支持修改配置
- [x] **语音识别：** 支持私聊接收和处理语音消息，通过文字回复

# 使用效果
## 私聊
![image](https://user-images.githubusercontent.com/66697106/232195996-fd5cfd40-82ab-4329-95c5-ae828762cba6.png)

## 群聊
![image](https://user-images.githubusercontent.com/66697106/232195808-1b2acfe4-01bd-4c79-9ce4-7ca2d2a67da4.png)

# 快速开始
支持 Linux、MacOS、Windows 系统（可在Linux服务器上长期运行),不需安装安装任何环境,如果是本地代码运行，需要安装golang

```shell
# 获取项目
git clone https://github.com/itxiaolin/openai-wechat.git
# 进入项目目录
cd openai-wechat
# 修改配置(配置api_key)
open-ai:
    api-key: 你的api_key
    base-url: https://api.openai.com/v1
# 依赖下载
go mod tidy 
# 启动项目
go run main.go
```

## 默认配置
```yaml
#只需要改open-ai，其他配置可按默认
open-ai:
    api-key: 你的api_key
    base-url: https://api.openai.com/v1

wx-robot:
    auto-pass: true
    session-timeout: 900
    storage-path: config/storage.json
    retry-num: 3
    context-cache-num: 10
    chatGPT-model: gpt-3.5-turbo
    voice:
        voice-dir: config/voice
    keyword-prefix:
        image-prompt: /image-prompt
        system-role: /system-role
        system-help: /system-help
        reset-context: /reset-context
```


## 命令使用教程
-  metisAI-wechat start  启动应用
-  metisAI-wechat start -d 后台启动应用
-  metisAI-wechat --config config/config.yaml start -d 指定配置文件启动
-  metisAI-wechat status 查看应用启动情况
-  metisAI-wechat stop 关闭应用
-  metisAI-wechat restart 重启应用

## 联系
欢迎提交PR、Issues，以及Star支持一下。如果你想了解更多项目细节，并与开发者们交流更多关于AI技术的实践，请加我好友。
![image](https://user-images.githubusercontent.com/66697106/232264989-f6adf8ee-e7cc-4cba-8afb-6d90c3c036cb.png)

