package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/cheggaaa/pb/v3"
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
		numTorrents := len(args)
		wg.Add(numTorrents)

		for _, arg := range args {
			go download(client, &wg, arg)
		}

		wg.Wait()
	},
}

func download(client *torrent.Client, wg *sync.WaitGroup, arg string) {
	defer wg.Done()

	if len(arg) > 6 && arg[:6] == "magnet" {
		tor, err := client.AddMagnet(arg)

		if err != nil {
			log.Fatal("invalid magnet")
			return
		}

		<-tor.GotInfo()
		tor.DownloadAll()

		done := make(chan bool)
		go client.WaitAll()
		go printProgress(tor, done)

		<-done
	} else {
		infoHash := fromInfoHashString(arg)
		tor, _ := client.AddTorrentInfoHash(infoHash)

		<-tor.GotInfo()
		tor.DownloadAll()

		done := make(chan bool)
		go client.WaitAll()
		go printProgress(tor, done)

		<-done
	}
}

func printProgress(tor *torrent.Torrent, done chan bool) {
	fmt.Println(tor.Name())
	length := tor.Length()

	reader := io.LimitReader(tor.NewReader(), length)
	writer := ioutil.Discard

	bar := pb.Simple.Start64(length)
	barReader := bar.NewProxyReader(reader)
	io.Copy(writer, barReader)

	bar.Finish()

	done <- true
}

func initClientConfig() *torrent.ClientConfig {
	home, _ := os.UserHomeDir()
	config := torrent.NewDefaultClientConfig()

	if runtime.GOOS == "windows" {
		config.DataDir = home + "\\Downloads"
	} else {
		config.DataDir = home + "/Downloads"
	}

	return config
}

func fromInfoHashString(hexString string) torrent.InfoHash {
	var infoHash torrent.InfoHash
	err := infoHash.FromHexString(hexString)

	if err != nil {
		panic("Invalid infohash")
	}

	return infoHash
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
