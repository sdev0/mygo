package kits

import (
	"os"
	"strings"

	"github.com/sdev0/mygo/sdk"
)

var (
	char_CH2EN = [2][]string{
		{"（", "）", "：", "；", "，", "‘", "’", "“", "”", "\n", " "},
		{"(", ")", ":", ";", ",", "'", "'", "\"", "\"", " ", " "},
	}
)

func FormatFiles(inpath string, outpath string, switch1_2 bool) {
	bytes, _ := os.ReadFile(inpath)
	lines := []string(strings.Split(string(bytes), "\n"))
	out := sdk.CreateFileByPath(outpath, os.O_CREATE|os.O_RDWR|os.O_APPEND)
	defer out.Close()
	var contentArr [][]string
	maxLen := 0
	for _, line := range lines {
		for i := range char_CH2EN[0] {
			line = strings.ReplaceAll(line, char_CH2EN[0][i], char_CH2EN[1][i])
		}
		for atmp := 0; atmp <= 5; atmp++ {
			line = strings.ReplaceAll(line, "  ", " ")
		}
		contentArr = append(contentArr, strings.SplitN(line, " ", 5))
		maxLen = Maxer(len(contentArr), maxLen)
	}
	var lenArr = make([]int, maxLen + 1)
	for _, arr := range contentArr {
		for j, str := range arr {
			lenArr[j] = Maxer(len(str), lenArr[j])
		}
	}
	// 是否进行第1，2项互换
	if switch1_2 {
		for i := range contentArr {
			if len(contentArr[i]) >= 2 {
				swapper(&contentArr[i][0], &contentArr[i][1])
			}
		}
		swapperInt(&lenArr[0], &lenArr[1])
	}
	for i := range contentArr {
		for j := range contentArr[i] {
			difLen := lenArr[j] - len(contentArr[i][j])
			AppenSpaceAfterString(&contentArr[i][j], difLen)
		}
	}
	for _, Arr := range contentArr {
		res := strings.Join(Arr, "  ")
		res = strings.TrimSpace(res)
		out.WriteString(res + "\n")
	}
	out.WriteString("\n")
}

func Maxer(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func swapper(str1, str2 *string) {
	tmp := *str1
	*str1 = *str2
	*str2 = tmp
}

func swapperInt(int1, int2 *int) {
	tmp := *int1
	*int1 = *int2
	*int2 = tmp
}

func AppenSpaceAfterString(str *string, size int) {
	if len(*str) == 0 {
		return
	}
	for ; size > 0; size-- {
		*str = *str + " "
	}
}
