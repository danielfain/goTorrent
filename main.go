package main

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/danielfain/goTorrent/cmd"
)

func main() {
	title := figure.NewFigure("goTorrent", "doom", true)
	title.Print()
	cmd.Execute()
}
