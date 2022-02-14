package kits

import "github.com/sdev0/mygo/sdk"

var ( // log
	mylog  *sdk.MyLog
	Log    func(v ...interface{})
	Logf   func(format string, v ...interface{})
	Linfo  func(v ...interface{})
	Linfof func(format string, v ...interface{})
	Lerr   func(v ...interface{})
	Lerrf  func(format string, v ...interface{})
)

func init() {
	mylog = sdk.InitLog("./log/", sdk.Ldate|sdk.Ltime|sdk.Lshortfile)
	Log = mylog.Llog
	Logf = mylog.Llogf
	Linfo = mylog.LogInfo
	Linfof = mylog.LogInfof
	Lerr = mylog.LogErr
	Lerrf = mylog.LogErrf

	mylog.SetShowHeader(false)
}

func DoKit() {
	// inpath, outpath := "./static/textFiles/sqlTable.txt", "./static/textFiles/sqlout.txt"
	// FormatFiles(inpath, outpath, true)
}
