//  @Description  : 使用主Config进行配置
//  @Author       : Shi
//  @Date         : 2022-01-22 16:29:39
//  @LastEditTime : 2022-01-22 16:29:39
//  @FilePath     : \mygo\config\config.go
package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sdev0/mygo/sdk"
)

type NovelConf struct {
	NovelWebsite        string   `toml:"NOVEL_WEBSITE"`
	NovelBaseUrl        string   `toml:"NOVEL_BASE_URL"`
	NovelDir            string   `toml:"NOVEL_DIR"`
	NovelName           []string `toml:"NOVEL_NAME"`
	NovelURL          []string `toml:"NOVEL_URL"`
	NovelResultJsonPath string   `toml:"NOVEL_RESULT_JSON_PATH"`
	Url_Append          string   `toml:"NOVEL_URL_CHAPTER_APPEND"`
	ThreadNum           int      `toml:"NOVEL_THREAD_NUM"`
	ChapterConstant     bool     `toml:"NOVEL_CHAPTER_CONSTANT=false"`
}
type SpiderConf struct {
	SpiderResDir string `toml:"SPIDER_RES_DIR"`
	SpiderImgDir string `toml:"SPIDER_IMAGE_DIR"`
}
type JavdiskConf struct {
	JavCoverUrl []string `toml:"JAV_COVER_URL"`
	ImgDir      string   `toml:"IMG_DIR"`
}

// 主Config
//  @param  [Novel]   novel配置
//  @param  [Spider]  spider配置
//  @param  [Javdisk] javdisk网站配置
type Config struct {
	Novel   NovelConf   `toml:"novel"`
	Spider  SpiderConf  `toml:"spider"`
	Javdisk JavdiskConf `toml:"javdisk"`
}

// 初始化自定义配置，配置必须是toml，且结构体必须含有初始化toml的信息
//  @param  confPath [string] 配置路径
//  @param  config [*interface{}] 配置结构体
//  @param  logger [*sdk.MyLog] 日志
func InitCustomConfig(confPath string, config interface{}, logger *sdk.MyLog) {
	_, err := toml.DecodeFile(confPath, config)
	if err != nil {
		logger.LogErr("Decode Config Failed.", err)
		os.Exit(1)
	}
}
