package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a torrent from a magnet or infohash",
	Run: func(cmd *cobra.Command, args []string) {
		config := initClientConfig()
		client, _ := torrent.NewClient(config)
		defer client.Close()

		var wg sync.WaitGroup
		wg.Add(len(args))

		for _, arg := range args {
			go download(client, &wg, arg)
		}

		wg.Wait()
	},
}

func download(client *torrent.Client, wg *sync.WaitGroup, arg string) {
	defer wg.Done()

	if arg[:6] == "magnet" {
		tor, err := client.AddMagnet(arg)

		if err != nil {
			log.Fatal("invalid magnet")
			return
		}

		<-tor.GotInfo()
		tor.DownloadAll()

		success := client.WaitAll()

		if !success {
			log.Fatal("error while downloading from infohash")
		}

	} else {
		infoHash := fromInfoHashString(arg)
		tor, _ := client.AddTorrentInfoHash(infoHash)

		<-tor.GotInfo()
		tor.DownloadAll()

		success := client.WaitAll()

		if !success {
			log.Fatal("error while downloading from infohash")
		}

	}

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

func init() {
	rootCmd.AddCommand(downloadCmd)
}
