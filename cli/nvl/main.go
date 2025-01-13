package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hkupty/neovim-bridge/share"
)

type Target string

const (
	Tab    Target = "t"
	Global Target = ""
)

func main() {
	client, err := share.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	location := "."
	if len(os.Args) > 1 {
		location = os.Args[len(os.Args)-1]
	}
	absoluteLocation, err := filepath.Abs(location)

	if err != nil {
		panic(err)
	}

	target := Tab

	if len(os.Args) > 2 && os.Args[1] == "global" {
		target = Global
	}

	var commandStr strings.Builder
	commandStr.WriteString(string(target))
	commandStr.WriteString("cd ")
	commandStr.WriteString(absoluteLocation)

	client.Exec(commandStr.String(), false)

	client.Wait()
	os.Exit(client.Exit)
}
