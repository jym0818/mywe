package ioc

import "github.com/jym0818/mywe/internal/service/oauth2/wechat"

func InitWechat() wechat.WechatService {
	service := wechat.NewwechatService("123456", "123456")
	return service
}
