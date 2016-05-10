# insco [![Build Status][travis-image]][travis-url]

insco is an installer for CLI developer tools, emacs, vim.
insco installs tools locally instead of system-wide.

## Usage

### Setup

insco installs tools to $HOME/bin. You need to append $HOME/bin to $PATH.

### Install tools

```shell
 $ insco vim
```

You can specify a version (WIP).

```shell
 $ insco vim 7.4
```

## Support

### Editors
- Emacs
- Vim

### Version Control Systems
- Git

### Tools
- [peco](https://github.com/peco/peco)
  - Interactive filtering tool
- [ghq](https://github.com/motemen/ghq)
  - Manage remote repository clones

[travis-image]: https://img.shields.io/travis/tatsuyafw/insco.svg
[travis-url]: https://travis-ci.org/tatsuyafw/insco
