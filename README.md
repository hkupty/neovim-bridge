# neovim-bridge

Simple $EDITOR bridge for neovim

---

neovim-bridge (nvb for short) is a very simple `$EDITOR` bridge for neovim.

It is a subset of [neovim-remote](https://github.com/mhinz/neovim-remote), aiming to address a simple yet common use-case, use neovim as `$EDITOR` for a process started from the `:terminal`, waiting until the file is discarded.

It has one fundamental difference over `neovim-remote`: the files are created setting `bufhidden=delete`, which means once you `:x` on that buffer, `nvb` will get response and release the process.
