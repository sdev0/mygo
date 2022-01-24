## NovelSpider.go

```go
func getChapterContent(aimurl string, file *os.File) error { 
	// Request the HTML page.
	res, err := http.Get(aimurl)
	if err != nil {
		Lerr(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		Lerrf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	// get title
	title := doc.Find(".entry-title").Text()
	file.WriteString(title + "\n")
	doc.Find(".entry-content").Find("p").Each(func (i int,s *goquery.Selection)  {
		file.WriteString(s.Text() + "\n")
	})
	return nil
}

import(
    "github.com/PuerkitoBio/goquery"
)

func HttpGet() string {
	ips := []string {
		"127.0.0.1:3500",
		"127.0.0.1:3520",
		"127.0.0.1:3585",
	}
	return ips[1]
}
func getXbookChapterByHttp(aimurl string) []Chapter {
	// Request the HTML page.
	res, err := http.Get(aimurl)
	if err != nil {
		Lerr(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		Lerrf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		Lerr(err)
	}

	var chapters []Chapter
	// Find the review items
	doc.Find(".entry-title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		url, _ := s.Find("a").Attr("href")
		chapters = append(chapters, Chapter{Title: title, Url: url})
		Log("Review %d: %s %s\n", i, title, url)
	})
	return chapters
}


func SpiderXbook() {
	pConf = Config.Novel
	initDownloadInfo(pConf.NovelResultJsonPath)
	for _, novelName := range pConf.NovelName {
		if _, ok := novelMap[novelName]; !ok {
			novelMap[novelName] = []Chapter{}
		}
		chapters := getXbookChapterByHttp(pConf.NovelBaseUrl + novelName)
		file := sdk.CreateFileByPath(pConf.NovelDir+novelName+".txt", os.O_CREATE|os.O_RDWR|os.O_APPEND)
		Linfof("正在获取 %s, 共 %d 章\n", novelName, len(chapters))
		for i := range chapters {
			if i > len(novelMap[novelName]) {
				if err := getChapterContentByHttp(chapters[i].Url, file); err == nil {
					Linfof("正在获取 %s, %s\n", novelName, chapters[i].Title)
					chs := novelMap[novelName]
					chs = append(chs, chapters[i])
					novelMap[novelName] = chs
				}
			}
		}
	}
	novelDownInfo = NovelDownloadInfo{}
	for name, chapters := range novelMap {
		novelDownInfo.NovelInfos = append(novelDownInfo.NovelInfos, NovelInfo{Name: name, Chapters: chapters})
	}
	bytes, err := json.Marshal(novelDownInfo)
	if err != nil {
		Lerr("配置marshal失败", err.Error())
	} else {
		file := sdk.CreateFileByPath(pConf.NovelResultJsonPath, os.O_CREATE|os.O_WRONLY)
		file.Write(bytes)
	}
}

cl.OnRequest(func(r *colly.Request) {
    r.Headers.Set("max-results", "1000")
})
```