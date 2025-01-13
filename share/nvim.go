package share

import (
	"fmt"
	"os"
	"sync"

	"github.com/neovim/go-client/nvim"
)

type NvimClient struct {
	*nvim.Nvim

	channel int
	wg      sync.WaitGroup
	Exit    int
}

func NewClient() (*NvimClient, error) {
	nvimAddress := os.Getenv("NVIM")
	cl, err := nvim.Dial(nvimAddress)

	if err != nil {
		return nil, err
	}

	client := NvimClient{Nvim: cl, wg: sync.WaitGroup{}, channel: cl.ChannelID()}

	return &client, nil
}

func (nvcl *NvimClient) RegisterBufferHandlers() error {
	err := nvcl.RegisterHandler("Exit", func(cl *nvim.Nvim, exitCode int) error {
		defer cl.Close()
		nvcl.wg.Done()

		fmt.Printf("Exit\n")

		nvcl.Exit = exitCode

		return nil
	})

	if err != nil {
		nvcl.Nvim.Close()
		return err
	}

	err = nvcl.RegisterHandler("Delete", func(cl *nvim.Nvim, modified bool) error {
		defer cl.Close()
		nvcl.wg.Done()

		if modified {
			nvcl.Exit = 1
		}

		return nil
	})

	if err != nil {
		nvcl.Nvim.Close()
		return err
	}

	return nil
}

func (nvcl *NvimClient) IncreaseLock() { nvcl.wg.Add(1) }
func (nvcl *NvimClient) Wait()         { nvcl.wg.Wait() }

func (nvcl *NvimClient) PrepareBuffer() (*nvim.Buffer, error) {
	buf, err := nvcl.CurrentBuffer()

	if err != nil {
		return nil, err
	}

	batch := nvcl.NewBatch()

	batch.SetBufferOption(buf, "bufhidden", "delete")
	batch.Command("augroup nvb")
	batch.Command(fmt.Sprintf("autocmd VimLeave * if exists(\"v:exiting\") && v:exiting > 0 | call rpcnotify(%d, \"Exit\", v:exiting) | endif", nvcl.channel))
	batch.Command(fmt.Sprintf("autocmd BufDelete <buffer=%d> silent! call rpcnotify(%d, \"Delete\", &modified)", buf, nvcl.channel))
	batch.Command("augroup END")

	if err = batch.Execute(); err != nil {
		return nil, err
	}

	return &buf, nil
}

func (nvcl *NvimClient) NewWindow() (*nvim.Window, error) {

	buf, err := nvcl.CreateBuffer(false, true)

	if err != nil {
		return nil, err
	}

	err = nvcl.Command("vspl")
	if err != nil {
		return nil, err
	}

	window, err := nvcl.CurrentWindow()

	if err != nil {
		return nil, err
	}

	err = nvcl.SetBufferToWindow(window, buf)

	if err != nil {
		return nil, err
	}

	return &window, nil
}
