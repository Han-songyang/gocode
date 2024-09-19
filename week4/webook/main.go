package main

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"log"
)

func main() {
	initViper()
	initLogger()
	server := InitWebServer()

	_ = server.Run(":8080")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	// 当前工作目录的 config 子目录
	viper.AddConfigPath("config")
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

//func initViperWatch() {
//	cfile := pflag.String("config",
//		"config/dev.yaml", "配置文件路径")
//	// 这一步之后，cfile 里面才有值
//	pflag.Parse()
//	//viper.Set("db.dsn", "localhost:3306")
//	// 所有的默认值放好s
//	viper.SetConfigType("yaml")
//	viper.SetConfigFile(*cfile)
//	viper.WatchConfig()
//	viper.OnConfigChange(func(in fsnotify.Event) {
//		log.Println(viper.GetString("test.key"))
//	})
//	// 读取配置
//	err := viper.ReadInConfig()
//	if err != nil {
//		panic(err)
//	}
//	val := viper.Get("test.key")
//	log.Println(val)
//}

//func initViperRemote() {
//	err := viper.AddRemoteProvider("etcd3",
//		"http://127.0.0.1:12379", "/webook")
//	if err != nil {
//		panic(err)
//	}
//	viper.SetConfigType("yaml")
//	viper.OnConfigChange(func(in fsnotify.Event) {
//		log.Println("远程配置中心发生变更")
//	})
//	go func() {
//		for {
//			err = viper.WatchRemoteConfig()
//			if err != nil {
//				panic(err)
//			}
//			log.Println("watch", viper.GetString("test.key"))
//			//time.Sleep(time.Second)
//		}
//	}()
//	err = viper.ReadRemoteConfig()
//	if err != nil {
//		panic(err)
//	}
//}
