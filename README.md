# neovim-bridge

Simple $EDITOR bridge for neovim

---

neovim-bridge (nvb for short) is a very simple `$EDITOR` bridge for neovim.

It is a subset of [neovim-remote](https://github.com/mhinz/neovim-remote), aiming to address a simple yet common use-case, use neovim as `$EDITOR` for a process started from the `:terminal`, waiting until the file is closed.

It has one fundamental difference over `neovim-remote`: the files are created setting `bufhidden=delete`, which means once you `:x` on that buffer, `nvb` will get a notification from neovim and release the process.

It also can be used on terminals outside neovim, spawning a neovim instane in place

### Usage

`nvb <file>` - open file in a new vertical tab; If neovim is not open, runs neovim and opens the file;

`<command> | nvb` - Opens the output of command in neovim;

`EDITOR=nvb <command>` - runs command using `nvb` as editor, for example `EDITOR=nvb git commit`.

`nvl` - sets current shell directory as the tab directory in neovim

`nvl <path>` - sets path as the tab directory in neovim
