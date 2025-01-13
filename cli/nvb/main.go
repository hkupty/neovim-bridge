package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hkupty/neovim-bridge/share"
)

var (
	isPipe   bool
	filename string
)

func init() {
	i, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	isPipe = (i.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func main() {
	client, err := share.NewClient()

	if err != nil {
		panic(err)
	}

	if err = client.RegisterBufferHandlers(); err != nil {
		panic(err)
	}

	defer client.Close()

	if isPipe {
		file, err := os.CreateTemp("", "nvb-*")
		if err != nil {
			panic(err)
		}

		filename = file.Name()
		defer os.Remove(filename)
		_, err = io.Copy(file, os.Stdin)

		if err != nil {
			panic(err)
		}
	} else {
		if len(os.Args) > 1 {
			filename = os.Args[1]
		} else {
			panic("No file")
		}
	}

	_, err = client.NewWindow()

	if err != nil {
		panic(err)
	}

	if filename == "" {
		panic(fmt.Sprintf("Filename is %s, isPipe is %t", filename, isPipe))
	}

	_, err = client.Exec("edit "+filename, false)
	client.IncreaseLock()

	if err != nil {
		panic(err)
	}

	_, err = client.PrepareBuffer()

	if err != nil {
		panic(err)
	}

	client.Wait()
	os.Exit(client.Exit)
}
