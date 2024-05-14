package main

import (
	"io"
	"os"
	"sync"
)

var (
	isPipe      bool
	nvimAddress string
	thisId      int
	wg          sync.WaitGroup
	file        *os.File
	exit        int
)

func init() {
	i, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	wg = sync.WaitGroup{}
	isPipe = (i.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func main() {
	nvimAddress = os.Getenv("NVIM")
	client, err := newClient()

	if err != nil {
		panic(err)
	}
	defer client.Close()

	thisId = client.ChannelID()

	if isPipe {
		file, err = os.CreateTemp("", "nvb-*")
		if err != nil {
			panic(err)
		}
		defer os.Remove(file.Name())
		_, err = io.Copy(file, os.Stdin)

		if err != nil {
			panic(err)
		}
	} else {
		if len(os.Args) > 1 {
			file, err = os.OpenFile(os.Args[1], os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
		} else {
			panic("No file")
		}
	}

	_, err = newWindow(client)

	if err != nil {
		panic(err)
	}

	_, err = client.Exec("edit "+file.Name(), false)
	wg.Add(1)

	if err != nil {
		panic(err)
	}

	_, err = prepareBuffer(client)

	if err != nil {
		panic(err)
	}

	wg.Wait()
	os.Exit(exit)
}
