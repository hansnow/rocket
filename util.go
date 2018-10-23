package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	SiteID string // PT站点上的种子ID
	URL    string // 种子链接
	Expire string // 过期时间
}

// RssTorrent Rss返回的种子信息
type RssTorrent struct {
	SiteID string // PT站点上的种子ID
	URL    string // 种子链接
	Size   int64  // 种子大小(Bytes)
}

type rssItem struct {
	Link string `xml:"url,attr"`
	Size int64  `xml:"length,attr"`
}
type rssResponse struct {
	Items []rssItem `xml:"channel>item>enclosure"`
}

// GetFreeTorrents 获取免费种子列表
// 具体流程分为三个步骤
// 1. 从PT站种子页面获取该页面所有种子，并筛选出Free的种子
// 2. 从RSS链接中获取输出的所有种子和对应的文件大小
// 3. 将所有Free种子中能在RSS中匹配到文件大小的种子提取出来并返回
func GetFreeTorrents() []FreeTorrent {
	freeTorrents := []FreeTorrent{}
	websiteTorrents := GetWebsiteTorrents()
	rssTorrents := GetRssTorrents()
	// 以websiteTorrents为准，把所有能拿到Size的种子留下来
	for _, wt := range websiteTorrents {
		hasSize := false
		for _, rt := range rssTorrents {
			if wt.SiteID == rt.SiteID {
				hasSize = true
				break
			}
		}
		if hasSize == true {
			freeTorrents = append(freeTorrents, wt)
		}
	}
	return freeTorrents
}

// GetWebsiteTorrents 获取网页免费种子列表
func GetWebsiteTorrents() []FreeTorrent {
	websiteTorrents := []FreeTorrent{}
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
		expire, exists := torrent.Find("b>span").Attr("title")
		// expire, _ := time.Parse("2006-01-02 15:04:05", expireStr)
		if exists == true && expire != "" {
			// 有些不符合这个条件是因为只写了免费，没有写时限
			websiteTorrents = append(websiteTorrents, FreeTorrent{
				SiteID: id,
				URL:    fmt.Sprintf("https://ourbits.club/download.php?id=%s&passkey=%s&https=1", id, passkey),
				Expire: expire,
			})
		}
	})
	return websiteTorrents
}

// GetRssTorrents 获取Rss种子列表
func GetRssTorrents() []RssTorrent {
	rssTorrents := []RssTorrent{}
	site := viper.GetStringMapString("site")
	passkey := site["passkey"]
	c := req.New()
	rssURL := fmt.Sprintf("https://ourbits.club/torrentrss.php?rows=50&passkey=%s&https=1&linktype=dl", passkey)
	r, _ := c.Get(rssURL)
	rssResp := rssResponse{}
	r.ToXML(&rssResp)
	for _, item := range rssResp.Items {
		re, _ := regexp.Compile("id=(\\d+)&passkey")
		siteID := re.FindStringSubmatch(item.Link)[1]
		rssTorrents = append(rssTorrents, RssTorrent{
			SiteID: siteID,
			Size:   item.Size,
		})
	}
	return rssTorrents
}

// TorrentRecord 种子记录
// records.csv格式
// `ID, SiteID, HashString, DateCreated, DateFreeEnd`
type TorrentRecord struct {
	ID          int64     // 种子ID，从1开始自增
	SiteID      string    // PT站点上的种子ID
	HashString  string    // 种子hash
	DateCreated time.Time // 添加时间
	DateFreeEnd time.Time // 免费结束时间
}

func getRecordsPath() string {
	// exePath, _ := os.Executable()
	// exeDir := filepath.Dir(exePath)
	// return filepath.Join(exeDir, "records.csv")
	return "records.csv"
}

// ReadTorrentRecords 读取rocket添加的种子记录
func ReadTorrentRecords() []TorrentRecord {
	result := []TorrentRecord{}
	recordsPath := getRecordsPath()
	// 文件不存在则创建文件
	f, err := os.OpenFile(recordsPath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	for _, r := range records {
		id, _ := strconv.ParseInt(r[0], 10, 64)
		dateCreated, _ := time.Parse("2006-01-02 15:04:05", r[3])
		dateFreeEnd, _ := time.Parse("2006-01-02 15:04:05", r[4])
		result = append(result, TorrentRecord{
			ID:          id,
			SiteID:      r[1],
			HashString:  r[2],
			DateCreated: dateCreated,
			DateFreeEnd: dateFreeEnd,
		})
	}
	return result
}

// WriteTorrentRecords 写入种子记录
func WriteTorrentRecords(records []TorrentRecord) {
	recordStr := [][]string{}
	for _, r := range records {
		recordStr = append(recordStr, []string{
			strconv.FormatInt(r.ID, 10),
			r.SiteID,
			r.HashString,
			r.DateCreated.Format("2006-01-02 15:04:05"),
			r.DateFreeEnd.Format("2006-01-02 15:04:05"),
		})
	}
	recordsPath := getRecordsPath()
	f, err := os.Create(recordsPath)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	w.WriteAll(recordStr)
}
