package wxrobot

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/itxiaolin/openai-wechat/internal/application/wxrobot/handlers"
	"github.com/itxiaolin/openai-wechat/internal/core/application/worker"
	"github.com/itxiaolin/openai-wechat/internal/core/logger"
	"github.com/itxiaolin/openai-wechat/internal/global"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"runtime"
	"time"
)

func NewWXBotEngine() worker.Engine {
	return &robotEngine{}
}

type robotEngine struct {
	bot *openwechat.Bot
}

func (r *robotEngine) GracefullyShutdown() {
	r.bot.Exit()
}

func (r *robotEngine) Process() {
	r.bot = r.DefaultBot(openwechat.Desktop)
	if global.Config.WxRobot.StoragePath != "" &&
		r.bot.HotLogin(openwechat.NewFileHotReloadStorage(global.Config.WxRobot.StoragePath)) == nil {
		user, _ := r.bot.GetCurrentUser()
		handlers.BotLoginTimeMap[r.bot.UUID()] = time.Now().Unix()
		logger.Info(nil, "热登录成功", zap.Any("用户名", user.NickName))
	} else {
		if err := r.bot.Login(); err != nil {
			logger.Error(nil, "登录异常", zap.Error(err))
			return
		}
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	go func() {
		_ = r.bot.Block()
	}()
}

func (r *robotEngine) DefaultBot(prepares ...openwechat.BotPreparer) *openwechat.Bot {
	bot := openwechat.NewBot(context.Background())
	// 注册登陆二维码回调
	bot.UUIDCallback = r.QrCodeCallBack
	// 扫码回调
	bot.ScanCallBack = func(_ openwechat.CheckLoginResponse) {
		logger.Info(nil, "扫码成功,请在手机上确认登录")
	}
	// 登录回调
	bot.LoginCallBack = func(_ openwechat.CheckLoginResponse) {
		handlers.BotLoginTimeMap[bot.UUID()] = time.Now().Unix()
		logger.Info(nil, "登录成功")
	}
	// 心跳回调函数,默认的行为打印SyncCheckResponse
	bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
		logger.Info(nil, "心跳函数，", zap.String("uuid", bot.UUID()), zap.String("RetCode", resp.RetCode))
	}
	for _, prepare := range prepares {
		prepare.Prepare(bot)
	}
	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	return bot
}

// QrCodeCallBack 登录扫码回调
func (r *robotEngine) QrCodeCallBack(uuid string) {
	if runtime.GOOS != "linux" {
		openwechat.PrintlnQrcodeUrl(uuid)
	} else {
		q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
		fmt.Println(q.ToString(true))
	}
	qrcodeUrl := openwechat.GetQrcodeUrl(uuid)
	logger.Info(nil, "访问下面网址扫描二维码登录", zap.String("登录地址", qrcodeUrl))
}
