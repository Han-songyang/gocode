package domain

import "time"

type User struct {
	Id         int64
	Email      string
	Password   string
	Nickname   string
	Birthday   time.Time
	AboutMe    string
	WechatInfo WechatInfo
	Phone      string
	Ctime      time.Time // UTC 0 的时区
}
