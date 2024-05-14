package main

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

func newClient() (*nvim.Nvim, error) {
	cl, err := nvim.Dial(nvimAddress)

	if err != nil {
		return nil, err
	}

	err = cl.RegisterHandler("Exit", func(cl *nvim.Nvim) error {
		defer cl.Close()
		fmt.Println("Exit")
		wg.Done()

		return nil
	})

	if err != nil {
		cl.Close()
		return nil, err
	}

	err = cl.RegisterHandler("Delete", func(cl *nvim.Nvim, buf int) error {
		defer cl.Close()
		fmt.Println("Done")
		wg.Done()

		return nil
	})

	if err != nil {
		cl.Close()
		return nil, err
	}

	return cl, nil
}

func prepareBuffer(client *nvim.Nvim) (*nvim.Buffer, error) {
	buf, err := client.CurrentBuffer()

	if err != nil {
		return nil, err
	}

	batch := client.NewBatch()

	batch.SetBufferOption(buf, "bufhidden", "delete")
	batch.Command("augroup nvb")
	batch.Command(fmt.Sprintf("autocmd VimLeave * if exists(\"v:exiting\") && v:exiting > 0 | call rpcnotify(%d, \"Exit\", v:exiting) | endif", thisId))
	batch.Command(fmt.Sprintf("autocmd BufDelete <buffer=%d> silent! call rpcnotify(%d, \"Delete\", %d)", buf, thisId, buf))
	batch.Command("augroup END")

	if err = batch.Execute(); err != nil {
		return nil, err
	}

	return &buf, nil
}

func newWindow(client *nvim.Nvim) (*nvim.Window, error) {

	buf, err := client.CreateBuffer(false, true)

	if err != nil {
		return nil, err
	}

	err = client.Command("vspl")
	if err != nil {
		return nil, err
	}

	window, err := client.CurrentWindow()

	if err != nil {
		return nil, err
	}

	err = client.SetBufferToWindow(window, buf)

	if err != nil {
		return nil, err
	}

	return &window, nil

}
