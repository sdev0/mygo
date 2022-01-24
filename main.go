//go:build ignore
// +build ignore

package main

import (
	"github.com/sdev0/mygo/sdk"
)

var ( // log
	mylog  *sdk.MyLog                            = sdk.InitLog("./log/", sdk.Ldate|sdk.Ltime|sdk.Lshortfile)
	Log    func(v ...interface{})                = mylog.Llog
	Logf   func(format string, v ...interface{}) = mylog.Llogf
	Linfo  func(v ...interface{})                = mylog.LogInfo
	Linfof func(format string, v ...interface{}) = mylog.LogInfof
	Lerr   func(v ...interface{})                = mylog.LogErr
	Lerrf  func(format string, v ...interface{}) = mylog.LogErrf
)

func main() {
	Log("1")
}

//////////////// test ping ////////////////
