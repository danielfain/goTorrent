package cmd

import (
	"fmt"
	"os"

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

		if args[0][:6] == "magnet" {
			tor, _ := client.AddMagnet(args[0])
			<-tor.GotInfo()
			tor.DownloadAll()

			success := client.WaitAll()

			if success {
				fmt.Println("torrent successfully downloaded")
				return
			}

			panic("error while downloading from magnet")
		} else {
			infoHash := fromInfoHashString(args[0])
			tor, _ := client.AddTorrentInfoHash(infoHash)

			<-tor.GotInfo()
			tor.DownloadAll()

			success := client.WaitAll()

			if success {
				fmt.Println("torrent successfully downloaded")
				return
			}

			panic("error while downloading from infohash")
		}
	},
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
