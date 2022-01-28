package mycmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	_ "time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sdev0/mygo/config"
	"github.com/sdev0/mygo/sdk"
)

var pConf config.NovelConf
var novelDownInfo NovelDownloadInfo
var novelMap = make(map[string][]Chapter)
var novelChan = make(chan NovelInfo)
var spiderCnt int = 0

type NovelDownloadInfo struct {
	NovelInfos []NovelInfo `json:"novelInfos"`
}

type NovelInfo struct {
	Name     string    `json:"name"`
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

//////////// init region ////////////

func initConfig() {
	idx := Config.NovelAll.NovelIndex
	pConf = config.NovelConf{
		NovelWebsite:        Config.NovelAll.NovelWebsite[idx],
		NovelBaseUrl:        Config.NovelAll.NovelBaseUrl[idx],
		NovelDir:            Config.NovelAll.NovelDir[idx],
		NovelName:           Config.NovelAll.NovelName,
		NovelURL:            Config.NovelAll.NovelURL,
		NovelResultJsonPath: Config.NovelAll.NovelResultJsonPath[idx],
		Url_Append:          Config.NovelAll.Url_Append,
		ThreadNum:           Config.NovelAll.ThreadNum,
		ChapterConstant:     Config.NovelAll.ChapterConstant,
	}
}

func initDownloadInfo(confPath string) {
	novelDownInfo = NovelDownloadInfo{}
	if !sdk.PathExist(confPath) {
		return
	}
	bytes, err := os.ReadFile(confPath)
	if err != nil {
		Lerr("Read config file failed", err.Error())
		os.Exit(1)
	}
	json.Unmarshal(bytes, &novelDownInfo)
	for _, infos := range novelDownInfo.NovelInfos {
		novelMap[infos.Name] = infos.Chapters
	}
}

//////////// get novel channel ////////////

func getNovelByChannelWithNumber(getChapterContent func(string, int, *os.File) error) {
	for {
		res := <-novelChan
		novelName := res.Name
		chapters := res.Chapters
		file := sdk.CreateFileByPath(pConf.NovelDir+novelName+".txt", os.O_CREATE|os.O_RDWR|os.O_APPEND)
		Linfof("正在获取 %s, 共 %d 章\n", novelName, len(chapters))
		for i := range chapters {
			if i > len(novelMap[novelName]) {
				if err := getChapterContent(chapters[i].Url, i, file); err == nil {
					Linfof("正在获取%03d: %s, %s\n", i, novelName, chapters[i].Title)
					novelMap[novelName] = append(novelMap[novelName], chapters[i])
				} else {
					Lerr(novelName, err.Error())
					if pConf.ChapterConstant {
						break
					}
				}
			}
		}
		spiderCnt--
	}
}

func getNovelByChannel(getChapterContent func(string, *os.File) error) {
	for {
		res := <-novelChan
		novelName := res.Name
		chapters := res.Chapters
		file := sdk.CreateFileByPath(pConf.NovelDir+novelName+".txt", os.O_CREATE|os.O_RDWR|os.O_APPEND)
		Linfof("正在获取 %s, 共 %d 章\n", novelName, len(chapters))
		for i := range chapters {
			if i > len(novelMap[novelName]) {
				if err := getChapterContent(chapters[i].Url, file); err == nil {
					Linfof("正在获取: %s, %s\n", novelName, chapters[i].Title)
					novelMap[novelName] = append(novelMap[novelName], chapters[i])
				} else {
					Lerr(novelName, err.Error())
					if pConf.ChapterConstant {
						break
					}
				}
			}
		}
		spiderCnt--
	}
}

//////////// Spider Book ////////////
func spiderBook(getChapter func(string) ([]Chapter, error), byNumber bool, getChapterContent func(string, *os.File) error, getChapterContentWithNumber func(string, int, *os.File) error) error {
	initConfig()
	initDownloadInfo(pConf.NovelResultJsonPath)
	Linfo("start to spider...")
	for i := 1; i <= pConf.ThreadNum; i++ {
		if byNumber {
			go getNovelByChannelWithNumber(getChapterContentWithNumber)
		} else {
			go getNovelByChannel(getChapterContent)
		}
	}
	for _, novelName := range pConf.NovelName {
		if _, ok := novelMap[novelName]; !ok {
			novelMap[novelName] = []Chapter{}
		}
		chapters, err := getChapter(pConf.NovelBaseUrl + novelName + pConf.Url_Append)
		if err != nil {
			return err
		}
		spiderCnt++
		novelChan <- NovelInfo{Name: novelName, Chapters: chapters}
	}
	for {
		if spiderCnt <= 0 {
			break
		}
	}
	novelDownInfo = NovelDownloadInfo{}
	for name, chapters := range novelMap {
		novelDownInfo.NovelInfos = append(novelDownInfo.NovelInfos, NovelInfo{Name: name, Chapters: chapters})
	}
	bytes, err := json.Marshal(novelDownInfo)
	if err != nil {
		Lerr("配置marshal失败", err.Error())
		return err
	} else {
		file := sdk.CreateFileByPath(pConf.NovelResultJsonPath, os.O_CREATE|os.O_WRONLY)
		file.Write(bytes)
	}
	return nil
}

//////////// XBOOK ////////////

func getXbookChapter(url string) ([]Chapter, error) {
	cl := colly.NewCollector()
	var chapters []Chapter
	cl.OnHTML(".date-outer>.entry-title", func(h *colly.HTMLElement) {
		chapters = append(chapters, Chapter{Title: h.ChildText("a"), Url: h.ChildAttr("a", "href")})
	})
	var Err error = nil
	cl.OnError(func(_ *colly.Response, err error) {
		Lerr(err)
		Err = err
	})
	cl.Visit(url)
	return chapters, Err
}

func getXbookChapterContent(aimurl string, file *os.File) {
	cl := colly.NewCollector()
	title := ""
	cl.OnHTML(".entry-title", func(h *colly.HTMLElement) {
		title = h.Text
		file.WriteString(title + "\n")
	})
	cl.OnHTML(".entry-content", func(h *colly.HTMLElement) {
		h.ForEach("p", func(i int, s *colly.HTMLElement) {
			file.WriteString(s.Text + "\n")
		})
	})
	cl.OnError(func(_ *colly.Response, err error) {
		Lerr(err)
	})
	cl.Visit(aimurl)
}

func getXbookChapterByHttp(aimurl string) ([]Chapter, error) {
	// Request the HTML page.
	res, err := http.Get(aimurl)
	if err != nil {
		Lerr(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err := fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		Lerr(err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		Lerr(err)
		return nil, err
	}

	var chapters []Chapter
	// Find the review items
	doc.Find(".entry-title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		url, _ := s.Find("a").Attr("href")
		chapters = append(chapters, Chapter{Title: title, Url: url})
	})
	return chapters, nil
}

func getXBookChapterContentByHttp(aimurl string, file *os.File) error {
	// Request the HTML page.
	res, err := http.Get(aimurl)

	if err != nil {
		Lerr(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// Lerrf("status code error: %d %s", res.StatusCode, res.Status)
		return errors.New("fail to get " + aimurl)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	// get title
	title := doc.Find(".entry-title").Text()
	file.WriteString(title + "\n")
	doc.Find(".entry-content").Find("p").Each(func(i int, s *goquery.Selection) {
		file.WriteString(s.Text() + "\n")
	})
	return nil
}

func getXBookChapterContentByHttpWithNumber(aimurl string, number int, file *os.File) error {
	// Request the HTML page.
	res, err := http.Get(aimurl)

	if err != nil {
		Lerr(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// Lerrf("status code error: %d %s", res.StatusCode, res.Status)
		return errors.New("fail to get " + aimurl)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	// get title
	title := doc.Find(".entry-title").Text()
	if title[0] == '\n' {
		title = title[1:]
	}
	file.WriteString(fmt.Sprintf("第%d章 %s\n", number, title))
	doc.Find(".entry-content").Find("p").Each(func(i int, s *goquery.Selection) {
		file.WriteString(s.Text() + "\n")
	})
	return nil
}

func SpiderXbook() {
	err := spiderBook(getXbookChapter, true, getXBookChapterContentByHttp, getXBookChapterContentByHttpWithNumber)
	if err != nil {
		Lerr("获取小说失败")
	}
}

//////////// 92qb ////////////

func Spider92qb() {
	err := spiderBook(get92qbChapter, false, get92qbChapterContent, nil)
	if err != nil {
		Lerr("获取小说失败")
	}
}

func get92qbChapter(aimurl string) ([]Chapter, error) {
	cl := colly.NewCollector()
	cl.DetectCharset = true
	Linfof("get the url chapters: %s\n", aimurl)
	var chapters []Chapter
	cl.OnHTML(".mulu_list>li", func(h *colly.HTMLElement) {
		chapters = append(chapters, Chapter{Title: h.ChildText("a"), Url: h.ChildAttr("a", "href")})
	})
	var Err error = nil
	cl.OnError(func(_ *colly.Response, err error) {
		Lerr(err)
		Err = err
	})
	cl.Visit(aimurl)
	return chapters, Err
}

func get92qbChapterByHttp(aimurl string) ([]Chapter, error) {
	// Request the HTML page.
	res, err := http.Get(aimurl)
	if err != nil {
		Lerr(err)
		return nil, err
	}
	defer res.Body.Close()
	Linfof("get the url chapters: %s\n", aimurl)
	if res.StatusCode != 200 {
		Lerrf("url: %s, status code %d, error: %s", aimurl, res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		Lerr(err)
		return nil, err
	}

	var chapters []Chapter
	// Find the review items
	doc.Find(".mulu_list").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		url, _ := s.Find("a").Attr("href")
		chapters = append(chapters, Chapter{Title: title, Url: url})
	})
	return chapters, nil
}

func get92qbChapterContent(aimurl string, file *os.File) error {
	cl := colly.NewCollector()
	cl.DetectCharset = true
	title := ""
	cl.OnHTML(".h1title>h1 ", func(h *colly.HTMLElement) {
		title = h.Text
		file.WriteString(title + "\n")
	})
	cl.OnHTML("#htmlContent", func(h *colly.HTMLElement) {
		content := h.Text
		content = strings.ReplaceAll(content, "        show_style();\n         show_style2();\n", "")
		content = strings.ReplaceAll(content, "        \n", "")
		content = strings.ReplaceAll(content, "			上一页        返回目录        下一页\n", "")
		content = strings.ReplaceAll(content, "         show_style3();\n", "")
		content = strings.ReplaceAll(content, "    ", "\n\t")
		file.WriteString(content + "\n")
	})
	var er error
	cl.OnError(func(_ *colly.Response, err error) {
		Lerr(err)
		er = err
	})
	cl.Visit(aimurl)
	return er
}
