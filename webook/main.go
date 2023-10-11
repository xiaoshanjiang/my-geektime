package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	// 注意，要在 Goland 里面把对应的 work director 设置到 webook
	// 要把配置初始化放在最前面
	initViperV2Watch()
	server := InitWebServer()
	// 注册路由
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})
	server.Run(":8080")
}

func initViper() {
	viper.SetDefault("db.dsn",
		"root:root@tcp(localhost:3306)/mysql")
	// 读取的文件名字叫做 dev
	viper.SetConfigName("dev")
	// 读取的类型是 yaml 文件
	viper.SetConfigType("yaml")
	// 在当前目录的 config 子目录下
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV1() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperReader() {
	cfg := `
db:
dsn: "root:root@tcp(localhost:13316)/webook"

redis:
addr: "localhost:6379"
`
	// 读取的类型是 yaml 文件
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func initViperV2() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV2Watch() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV3Remote() {
	err := viper.AddRemoteProvider("etcd3",
		"http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			err = viper.WatchRemoteConfig()
			if err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(time.Second) // 睡个一秒钟
		}
	}()
}
