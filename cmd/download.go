package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a torrent from a magnet or infohash",
	Long:  "Download a torrent from a magnet or infohash",
	Run: func(cmd *cobra.Command, args []string) {
		t, _ := cmd.Flags().GetString("type")
		fmt.Println(t)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringP("type", "t", "", "Magnet or infohash")

}
