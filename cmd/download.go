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

		torrentType, _ := cmd.Flags().GetString("type")

		if torrentType == "infohash" || torrentType == "hash" {
			infoHash := fromInfoHashString(args[0])
			tor, _ := client.AddTorrentInfoHash(infoHash)
			<-tor.GotInfo()
			tor.DownloadAll()
			if client.WaitAll() == true {
				fmt.Println("torrent successfully downloaded")
			}
			fmt.Println("there was an error while downloading")
		}

		if torrentType == "magnet" {

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

	downloadCmd.Flags().StringP("type", "t", "", "Magnet or infohash")

}
