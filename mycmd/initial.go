package mycmd

import (
	"github.com/sdev0/mygo/config"
	"github.com/sdev0/mygo/sdk"
)

var (
	logger *sdk.MyLog     // 主日志
	Config *config.Config // 主Config
)

var (
	Log    func(v ...interface{})
	Logf   func(format string, v ...interface{})
	Linfo  func(v ...interface{})
	Linfof func(format string, v ...interface{})
	Lerr   func(v ...interface{})
	Lerrf  func(format string, v ...interface{})
)

func InitAll() {
	initLog()
	Config = new(config.Config)
	config.InitCustomConfig("./config/config.toml", Config, logger)
	//Log(Config.NovelConf, Config.JavdiskConf, Config.SpiderConf)
}
func initLog() {
	logpath := "./log/"
	logger = sdk.InitLog(logpath, sdk.Ldate|sdk.Ltime|sdk.Lshortfile)
	Log = logger.Llog
	Logf = logger.Llogf
	Linfo = logger.LogInfo
	Linfof = logger.LogInfof
	Lerr = logger.LogErr
	Lerrf = logger.LogErrf
}
