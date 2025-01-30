package main

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/hkupty/neovim-bridge/share"
)

// TODO use flag to handle command line flags

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

func prepareFileForPipe() (string, error) {
	file, err := os.CreateTemp("", "nvb-*")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, os.Stdin)

	if err != nil {
		os.Remove(file.Name())
		return "", err
	}

	return file.Name(), nil
}

func main() {
	var filename string
	var err error

	if isPipe {
		filename, err = prepareFileForPipe()
		if err != nil {
			panic(err) // TODO handler error correctly
		}
		defer os.Remove(filename)
	} else {
		if len(os.Args) > 1 {
			filename = os.Args[1]
		} else {
			slog.Warn("No file")
			os.Exit(4)
		}

		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			f, err := os.OpenFile(filename, os.O_CREATE, 0755)
			if err != nil {
				slog.Warn("Unable to create file", "error", err)
				os.Exit(2)
			}
			f.Close()
		}

	}

	client, err := share.NewClient()

	if err != nil {
		nvimPath, err := exec.LookPath("nvim")

		if err != nil {
			slog.Warn("Neovim doesn't seem to be installed, exiting", "error", err)
			os.Exit(3)
		}

		slog.Warn("Neovim might not be running. Spawning", "filename", filename)

		err = syscall.Exec(nvimPath, []string{"nvim", filename}, syscall.Environ())

		if err != nil {
			slog.Error("Neovim closed with error", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err = client.RegisterBufferHandlers(); err != nil {
		panic(err)
	}

	defer client.Close()

	_, err = client.NewWindow()

	if err != nil {
		panic(err)
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
