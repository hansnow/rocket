package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"github.com/spf13/viper"
)

// TrConfig Transmission配置
type TrConfig struct {
	URL      string `json:"url"`
	User     string
	Password string
}

// SiteConfig 站点配置
type SiteConfig struct {
	Name    string
	Passkey string
}

// Config 配置
type Config struct {
	Transmission TrConfig
	Site         SiteConfig
	SizeLimit    int `json:"size_limit"`
	StorageLimit int `json:"storage_limit"`
	Policy       string
}

// ReadConfig 读取配置
func ReadConfig() {
	exePath, _ := os.Executable()
	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Dir(exePath))
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// GetTorrents 获取所有种子列表

// FreeTorrent 站点免费种子信息
type FreeTorrent struct {
	URL    string // 种子链接
	Expire string // 过期时间
}

// GetFreeTorrents 获取免费种子列表
func GetFreeTorrents() []FreeTorrent {
	ret := []FreeTorrent{}
	// Get necessary value from config
	site := viper.GetStringMapString("site")
	cookie := site["cookie"]
	passkey := site["passkey"]
	c := req.New()
	cookieHeader := req.Header{"Cookie": cookie}
	r, _ := c.Get("http://ourbits.club/torrents.php", cookieHeader)
	doc, _ := goquery.NewDocumentFromResponse(r.Response())
	doc.Find("table.torrentname").Has("img.pro_free").Each(func(i int, torrent *goquery.Selection) {
		// 获取torrenet_id
		href, _ := torrent.Find("a[href^='download.php']").Attr("href")
		id := strings.Split(href, "id=")[1]
		// 获取免费剩余时间
		expire, _ := torrent.Find("b>span").Attr("title")
		// expire, _ := time.Parse("2006-01-02 15:04:05", expireStr)
		ret = append(ret, FreeTorrent{
			URL:    fmt.Sprintf("https://ourbits.club/download.php?id=%s&passkey=%s&https=1", id, passkey),
			Expire: expire,
		})
	})
	return ret
}
