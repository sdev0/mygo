package mycmd

import (
	"fmt"
	"os"

	"github.com/gocolly/colly"
	"github.com/sdev0/mygo/sdk"
)

type AimInfo struct {
	Title string
	Url   string
}

func SpiderJav() {
	//var urls []string
	var allInfos []AimInfo
	for i := 8; i <= 9; i++ {
		url := fmt.Sprintf("https://javdisk.com/studio/madonna/page-%d.html", i)
		res := getListCover(url)
		allInfos = append(allInfos, res...)
	}
	//Log(res)
	file := sdk.CreateFileByPath(Config.Spider.SpiderResDir+"JUL.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND)
	file.WriteString("-------------- " + sdk.GetTimeStrByKey(sdk.T_STR2) + " --------------\n")

	imgDir := Config.Spider.SpiderImgDir + Config.Javdisk.ImgDir
	sdk.CreateFileByPath(imgDir, 0)
	for _, pics := range allInfos {
		file.WriteString(pics.Title + "(" + pics.Url + ")\n")
		size, err := sdk.DownloadFile(imgDir, pics.Url, pics.Title, "", 0, 5)
		if err != nil {
			Lerr(fmt.Sprintf("Download %s failed: %s", pics.Title, err.Error()))
		} else {
			Logf("Download %s success, file size: %d", pics.Title, size)
		}
	}
}

// 获取javdisk网页列表中的封面
//  @param  aimurl [string] 目标URL
//  @return [[]AimInfo] 封面信息
func getListCover(aimurl string) []AimInfo {
	Linfo("正在爬取 " + aimurl)
	cl := colly.NewCollector()
	var aimInfos []AimInfo
	cl.OnHTML(".wrap-main-item", func(h *colly.HTMLElement) {
		url := h.ChildAttr("img", "data-src")
		title := h.ChildAttr("img", "alt")
		if len(title) > 180 {
			title = title[:180]
		}
		aimInfos = append(aimInfos, AimInfo{Title: title, Url: url})
	})
	cl.OnError(func(_ *colly.Response, err error) {
		Lerr(err)
	})
	cl.Visit(aimurl)
	defer Linfof("爬取内容共 %d 条", len(aimInfos))
	return aimInfos
}
