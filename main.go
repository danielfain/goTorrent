/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"log"
	"os"

	"github.com/anacrolix/torrent"
	"github.com/danielfain/goTorrent/cmd"
)

func main() {
	cmd.Execute()
	config := initClientConfig()
	client, _ := torrent.NewClient(config)
	defer client.Close()

	infoHash := fromInfoHashString("8543AC1F905954E0EC0F8E487646601C9F9F41CF")
	torrent, _ := client.AddTorrentInfoHash(infoHash)

	<-torrent.GotInfo()
	torrent.DownloadAll()
	client.WaitAll()
	log.Print("torrent finished downloading")
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
