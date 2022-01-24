//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/go-ping/ping"
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
	mylog.SetShowHeader(false)
	test2DArray()
}

//////////////// test ping ////////////////
func testPing() {
	ip := "172.16.2.3"
	res, err := PingConn_Linux(ip)
	if err != nil {
		Log("ip("+ip+"):", res)
	}
}

// 判断linux下能否ping ip成功
//  @param  addr [string]
//  @return [bool] 能够ping ip成功
//  @return [error] ping ip时出现的错误
func PingConn_Linux(ip string) (bool, error) {
	Command := fmt.Sprintf("ping -c 1 -W 3 %s > /dev/null && echo true || echo false", ip)
	output, err := exec.Command("/bin/sh", "-c", Command).Output()
	return string(output) == "true\n", err
}

// 测试linux下的ping接收信息
//  @param  ip [string]
func PingTest_Linux(ip string) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}
	pinger.Debug = true
	pinger.OnFinish = func(s *ping.Statistics) {
		Logf("OnFinish: %#v\n", s)
	}
	pinger.OnDuplicateRecv = func(p *ping.Packet) {
		Logf("OnDuplicateRecv: %#v\n", p)
	}
	pinger.OnRecv = func(p *ping.Packet) {
		Logf("OnRecv: %#v\n", p)
	}
	pinger.Timeout = time.Second * 3
	pinger.Count = 3
	pinger.Run()

}

//////////////// test two-dimensional array ////////////////
func test2DArray() {
	Log("####### test two-dimensional array")
	// define 2D array
	var twoDarr [2][]string
	if len(twoDarr[0]) != 0 || len(twoDarr[1]) != 0 {
		Linfo("2Darray is not null. length:", len(twoDarr[0]), len(twoDarr[1]))
	} else {
		Linfo("2Darray is null.")
	}
	Logf("len 2darray: %d, child array len: %d, %d, content of twoDarr: %+v", len(twoDarr), len(twoDarr[0]), len(twoDarr[1]), twoDarr)
	twoDarr[0] = append(twoDarr[0], "123")
	Logf("len 2darray: %d, content of twoDarr: %+v", len(twoDarr), twoDarr)
	twoDarr[1] = append(twoDarr[1], "456")
	Logf("len 2darray: %d, content of twoDarr: %+v", len(twoDarr), twoDarr)
}
