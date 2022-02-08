package test

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/go-ping/ping"
)

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
