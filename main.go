package main

import (
	"os"

	"github.com/anacrolix/torrent"
	"github.com/danielfain/goTorrent/cmd"
)

func main() {
	cmd.Execute()
	/*
		config := initClientConfig()
		client, _ := torrent.NewClient(config)
		defer client.Close()

		infoHash := fromInfoHashString("8543AC1F905954E0EC0F8E487646601C9F9F41CF")
		torrent, _ := client.AddTorrentInfoHash(infoHash)

		<-torrent.GotInfo()
		torrent.DownloadAll()
		client.WaitAll()
		log.Print("torrent finished downloading")
	*/
}

func initClientConfig() *torrent.ClientConfig {
	home, _ := os.UserHomeDir()
	config := torrent.NewDefaultClientConfig()
	config.DataDir = home + "\\Downloads"
	return config
}

func fromInfoHashString(hexString string) torrent.InfoHash {
	var infoHash torrent.InfoHash
	infoHash.FromHexString(hexString)
	return infoHash
}
