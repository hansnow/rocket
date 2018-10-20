package main

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ReadConfig()
	// To test transmissionrpc
	// tr := viper.GetStringMapString("transmission")
	// c, err := transmissionrpc.New(tr["url"], tr["user"], tr["password"], nil)
	// checkErr(err)
	// torrentList, err := c.TorrentGetAll()
	// checkErr(err)
	// for _, torrent := range torrentList {
	// 	fmt.Println(*torrent.Name)
	// }

	// To test goquery
	// t := GetFreeTorrents()
	// fmt.Println(t)

}
