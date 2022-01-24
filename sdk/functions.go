package sdk

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

/**
 * @description:
 *  判断当前路径是否存在
 * @param  path [string] 路径
 * @return bool 存在返回true，如果不存在返回false
 */
func PathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

/**
 * @description:
 *  处理ERROR信息
 * @param  err [error] 给定error
 * @param  exit [bool] error != nil时，是否退出程序
 * @param  content [...interface{}] error != nil时，输出的内容
 */
func CheckError(err error, exit bool, content ...interface{}) {
	if err != nil {
		if len(content) != 0 {
			fmt.Print(content...)
		}
		if exit {
			os.Exit(1)
		}
	}
}

/**
 * @description:
 *  执行命令command，并获取命令结果
 * @param  command [string] 待执行的命令
 * @param  printOut [bool] 执行命令时是否直接打印输出
 * @return [string] 命令执行结果
 * @return [error]  命令执行错误信息
 */
func GetCmdResult(command string, printOut bool) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	retOut := ""
	if err != nil {
		return retOut, err
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("命令启动失败")
		panic(err)
	}
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		retOut += line
		if printOut {
			fmt.Print(line)
		}
	}
	return retOut, nil
}

/**
* @description:
*  生成随机字符串数组
* @param  [*] 无
* @return [string] 生成的随机字符串
 */
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	rand.Seed(time.Now().Unix())
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 根据路径创建文件/目录，返回创建/打开的文件 or nil(目录)
//  @param  path [string] 如果是目录，需以 \ 或 / 结尾
//  @param  filemode [int] 
//  @return [*os.File] 
func CreateFileByPath(path string, filemode int) *os.File {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	CheckError(err, false, "[ERROR] 创建文件目录失败: ", dir)
	if ch := path[len(path) - 1]; ch == '\\' || ch == '/' {
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|filemode, 0666)
	CheckError(err, false, "[ERROR] 打开文件失败: ", path)
	return file
}

// 根据key获取时间字符
//  @param  key [int]
//	T_STR  = 1024 format:20060102150405
//	T_STR2 = 2048 format:2006-01-02 15:04:05
// 	T_MILISEC = 2    T_SEC = 4   T_MIN = 8    T_HOUR = 16
// 	T_TIME    = 32   T_DAY = 64  T_MON = 128  T_YEAR = 256
// 	T_DATE    = 512
//  @param  format [int]
//  format:0 20060102150405
//  format:1 2006-01-02 15:04:05
//  @return [*]
func GetTimeStrByKey(key int) string {
	timestr := ""
	timearr := [][2]string{
		{"", ""},                                       // 0
		{"", time.Now().Format(".999999999")},          // 1
		{"", time.Now().Format("05")},                  // 2
		{"", time.Now().Format("04")},                  // 3
		{"", time.Now().Format("15")},                  // 4
		{"", time.Now().Format("150405")},              // 5
		{"", time.Now().Format("02")},                  // 6
		{"", time.Now().Format("01")},                  // 7
		{"", time.Now().Format("2006")},                // 8
		{"", time.Now().Format("20060102")},            // 9
		{"", time.Now().Format("20060102150405")},      // 10
		{"", time.Now().Format("2006-01-02 15:04:05")}, // 11
	}
	timestr = timearr[11][(key>>11)&1] + timearr[10][(key>>10)&1] +
		timearr[9][(key>>9)&1] + timearr[8][(key>>8)&1] +
		timearr[7][(key>>7)&1] + timearr[6][(key>>6)&1] +
		timearr[5][(key>>5)&1] + timearr[4][(key>>4)&1] +
		timearr[3][(key>>3)&1] + timearr[2][(key>>2)&1] +
		timearr[1][(key>>1)&1] + timearr[0][(key)&1]
	return timestr
}

// time.Duration转换为00h00m00s
//  @param  dur [time.Duration]
//  @return [string] 返回格式为01h02m03s
func DurationToTimeString(dur time.Duration) string {
	second := int(dur.Seconds())
	minute := second / 60
	second = second % 60
	hour := minute / 60
	minute = minute % 60
	return fmt.Sprintf("%02dh%02dm%02ds", hour, minute, second)
}
