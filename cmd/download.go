package cmd

import (
	"io"
	"io/ioutil"
	"sync"

	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a torrent from a magnet or infohash",
	Run: func(cmd *cobra.Command, args []string) {
		config := initClientConfig()
		client, _ := torrent.NewClient(config)
		defer client.Close()

		numTorrents := len(args)

		var wg sync.WaitGroup
		torrents := make(chan *torrent.Torrent, numTorrents)

		wg.Add(numTorrents)

		for _, arg := range args {
			go getTorrentInfo(client, &wg, arg, torrents)
		}

		wg.Wait()

		p := mpb.New(mpb.WithWaitGroup(&wg))

		wg.Add(numTorrents)

		for t := range torrents {
			go printProgress(t, &wg, p)
		}

		wg.Wait()
		p.Wait()
	},
}

func getTorrentInfo(client *torrent.Client, wg *sync.WaitGroup, arg string, torrents chan *torrent.Torrent) {
	defer wg.Done()

	if len(arg) > 6 && arg[:6] == "magnet" {
		tor, err := client.AddMagnet(arg)

		if err != nil {
			panic("Invalid magnet.")
		}

		<-tor.GotInfo()
		torrents <- tor
	} else {
		infoHash := fromInfoHashString(arg)
		tor, _ := client.AddTorrentInfoHash(infoHash)

		<-tor.GotInfo()
		torrents <- tor
	}
}

func printProgress(tor *torrent.Torrent, wg *sync.WaitGroup, p *mpb.Progress) {
	defer wg.Done()

	total := tor.Length()
	name := tor.Name()

	bar := p.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(name+" "),
			decor.CountersKibiByte("% 6.1f / % 6.1f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.AverageSpeed(decor.UnitKiB, "% .2f"), "done",
			),
		),
	)

	torrentReader := io.LimitReader(tor.NewReader(), total)
	barReader := bar.ProxyReader(torrentReader)
	io.Copy(ioutil.Discard, barReader)
}

func initClientConfig() *torrent.ClientConfig {
	home, _ := homedir.Dir()
	config := torrent.NewDefaultClientConfig()
	config.DataDir = home + "/Downloads"
	config.Logger = log.Discard

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
