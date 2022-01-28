> 这是mygo的日志，用来记录每天干了什么

# 日志ing

## 2022.1.21

编写爬虫，主要爬取`javdisk.com`上的视频封面和视频名，并将图片以视频名下载下来

## 2022.1.22
打算：  
- [x] 完善javdisk爬虫，编写专门的函数
- [x] 完善`config.toml`读取函数
- [x] 完善xbook小说爬取
- [x] 完成xbook爬取小说章节记录

遇到的问题：  
- [x] 爬取网页并不会加载网页到网页尾部  
在网页后加上需要加载最大数量的信息。  
怎么找到的呢：开发者选项进行刷新，找到xbookcn提供的文件，当将网页拉下去的时候，将会显示新加载的文件，新加载的文件的文件名就是请求的链接，链接内包含最大请求章节数

## 2022.1.24

打算：  
- [x] 测试`ping`  
>测试均在linux下进行  
1. 使用linux下的`ping`连接成功命令
```
ping -c 1 -W 3 %s > /dev/null && echo true || echo false
```
2. 使用`go-ping/ping`接收ping信息
```go
pinger, err := ping.NewPinger(ip)                // 创建pinger
pinger.Debug = true                              // 设置是否需要调试
pinger.Timeout = time.Second                     // 设置单词ping时长
pinger.Count = 3                                 // 设置ping次数
pinger.OnFinish = func(s *ping.Statistics) {}    // ping结束后初始数据（多次ping结果）
pinger.OnDuplicateRecv = func(p *ping.Packet) {} // 处理已经被接收的数据包又被ping接收
pinger.OnRecv = func(p *ping.Packet) {}          // 处理单次ping接收数据包
pinger.Run()                                     // 执行ping
```

## 2022.1.26

日程：  
- [x] 完成layuimini部署在beego上
- [x] 完成92qb的小说爬虫
在使用colly爬取内容的时候，如果爬取下来的内容是乱码的，那么就在colly定义collector的下面，写上这么个语句：`colly.Collector.DetectCharset = true`

## 2022.1.28

日程  
- [x] 优化novel spider的配置和爬取函数



# 记录ing

## 教程
[gitignore语法](https://blog.csdn.net/qq_39109805/article/details/93379035)  
[colly使用文档](https://www.jianshu.com/p/cbe0f6aae5bf)