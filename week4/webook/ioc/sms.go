package ioc

import (
	"webook/internal/service/sms"
	"webook/internal/service/sms/localsms"
)

func InitSms() sms.Service {
	return localsms.NewService()
}
