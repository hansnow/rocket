package main

import (
	"flag"
	"fmt"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ReadConfig()
	// 使用 `-clear=true` 参数运行时，不做其他操作
	isClear := flag.Bool("clear", false, "clear all controlled torrents")
	flag.Parse()
	if *isClear == true {
		Clear()
		return
	}
	// To test transmissionrpc
	// tr := viper.GetStringMapString("transmission")
	// c, err := transmissionrpc.New(tr["url"], tr["user"], tr["password"], nil)
	// checkErr(err)
	// torrentList, err := c.TorrentGetAll()
	// checkErr(err)
	// for _, torrent := range torrentList[len(torrentList)-5:] {
	// 	fmt.Println(*torrent.ID, *torrent.Name, (*torrent.DateCreated).Format("2006-01-02 15:04:05"))
	// }

	// To test goquery
	// freeTorrents := GetFreeTorrents()
	// fmt.Println(len(freeTorrents))
	// for _, t := range freeTorrents {
	// 	fmt.Println(t.Expire, t.URL)
	// }

	// To test records.csv
	// ReadTorrentRecords()

	// To test add torrent
	// tr := viper.GetStringMapString("transmission")
	// c, err := transmissionrpc.New(tr["url"], tr["user"], tr["password"], &transmissionrpc.AdvancedConfig{
	// 	HTTPTimeout: time.Duration(60 * time.Second),
	// })
	// checkErr(err)
	// url := "<TORRENT_URL>"
	// downloadDir := tr["download-dir"]
	// paused := false
	// t, err := c.TorrentAdd(&transmissionrpc.TorrentAddPayload{
	// 	Filename:    &url,
	// 	DownloadDir: &downloadDir,
	// 	Paused:      &paused,
	// })
	// checkErr(err)
	// fmt.Println(*t.ID)
	// fmt.Println(*t.Name)
	// fmt.Println(*t.HashString)

	// To test rss torrent list
	ts := GetFreeTorrents()
	for _, t := range ts {
		fmt.Println(t.SiteID, t.URL, t.Expire)
	}
}
