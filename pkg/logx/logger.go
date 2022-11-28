package logx

import (
	"io"
	"supreme-flamego/config"
	"supreme-flamego/pkg/colorful"
	"log"
	"os"
)

type debugDefault struct {
	Debug *log.Logger
}

func (d *debugDefault) Println(v ...interface{}) {
	if config.GetConfig().MODE == "debug" {
		d.Debug.Println(v...)
	}
}

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *debugDefault
)

func InitLogger() {
	if config.GetConfig().LogPath == "" {
		log.Fatalln("LogPath 未设置")
	}
	if _, err := os.Stat(config.GetConfig().LogPath); os.IsNotExist(err) {
		log.Print("LogPath 不存在")
		if _, err := os.Create(config.GetConfig().LogPath); err != nil {
			log.Fatalln("新建 LogPath 失败 ", err)
		}
		log.Print("新建 LogPath 成功")
	}
	errFile, err := os.OpenFile(config.GetConfig().LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开 LogPath 失败 ", err)
	}

	Info = log.New(os.Stdout, "[Info] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, colorful.Yellow("[Warning] "), log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile), colorful.Red("[Error] "), log.Ldate|log.Ltime|log.Lshortfile)
	Debug = &debugDefault{
		Debug: log.New(os.Stdout, colorful.Blue("[Debug] "), log.Ldate|log.Ltime|log.Lshortfile),
	}
}
