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

# 记录ing

## 教程
[gitignore语法](https://blog.csdn.net/qq_39109805/article/details/93379035)  
[colly使用文档](https://www.jianshu.com/p/cbe0f6aae5bf)