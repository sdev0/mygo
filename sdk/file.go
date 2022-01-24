package sdk

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 下载文件
//  @param  fileurl [*] 文件下载连接 必要
//  @param  storedir [*] 存储路径 必要
//  @param  filename [*] 存储文件名 必要
//  @param  ext[string] 文件扩展名 如果为空，则根据url获取
//  @param  buffersize [int] 下载文件缓存大小，如果为0，则自定义为1024*1024
//  @param  try [int] 下载文件尝试次数
//  @return [int] 文件大小
//  @return [error] 错误信息
func DownloadFile(fileurl, storedir, filename, ext string, buffersize, try int) (int, error) {
	if ext == "" {
		ext = filepath.Ext(fileurl)
	}
	if PathExist(storedir + filename + ext) {
		return 0, nil
	}
	res, err := http.Get(fileurl)
	if err != nil {
		if try > 0 {
			try--
			DownloadFile(fileurl, storedir, filename, ext, buffersize, try)
		}
		return 0, err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	if buffersize <= 0 {
		buffersize = 1048576
	}
	reader := bufio.NewReaderSize(res.Body, buffersize)
	file, err := os.Create(storedir + filename + ext)
	if err != nil {
		return 0, err
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	fileLen, _ := io.Copy(writer, reader)
	return int(fileLen), nil
}
// 下载文件(简略版)，根据url获取文件名
//  @param  fileurl [*] 文件下载连接 必要
//  @param  storedir [*] 存储路径 必要
//  @param  buffersize [int] 下载文件缓存大小，如果为0，则自定义为1024*1024
//  @param  try [int] 下载文件尝试次数
//  @return [int] 文件大小
//  @return [error] 错误信息
func DownloadFileSimple(fileurl, storedir string, buffersize, try int) (int, error) {
	index := strings.LastIndex(fileurl, "/")
	if index == -1 {
		index = strings.LastIndex(fileurl, "\\")
	}
	filename := fileurl[index:]
	index = 1
	if filename == "" {
		filename = "unamedfile"
	}
	for PathExist(storedir + filename) {
		filename = fmt.Sprintf("unamedfile-%d", index)
	}
	res, err := http.Get(fileurl)
	if err != nil {
		if try > 0 {
			try--
			DownloadFileSimple(fileurl, storedir, buffersize, try)
		}
		return 0, err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	if buffersize <= 0 {
		buffersize = 1048576
	}
	reader := bufio.NewReaderSize(res.Body, buffersize)
	file, err := os.Create(storedir + filename)
	if err != nil {
		return 0, err
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	fileLen, _ := io.Copy(writer, reader)
	return int(fileLen), nil
}