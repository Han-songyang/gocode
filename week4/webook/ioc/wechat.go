package ioc

import (
	"webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	//appID, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("找不到环境变量 WECHAT_APP_ID")
	//}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("找不到环境变量 WECHAT_APP_SECRET")
	//}
	appID, appSecret := "test", "test"
	return wechat.NewService(appID, appSecret)
}
