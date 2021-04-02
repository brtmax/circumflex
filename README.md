<p align="center">
  <img src="screenshots/circumflex.png" width="350" alt="^"/>
</p>

#

<div align="center">

[![Latest release](https://img.shields.io/github/v/release/bensadeh/circumflex?label=stable&color=e1acff&labelColor=292D3E)](https://github.com/bensadeh/circumflex/releases)
[![Changelog](https://img.shields.io/badge/docs-changelog-9cc4ff?labelColor=292D3E)](https://github.com/bensadeh/circumflex/blob/master/CHANGELOG.md)
[![License](https://img.shields.io/github/license/bensadeh/circumflex?color=c3e88d&labelColor=292D3E)](https://github.com/bensadeh/circumflex/blob/master/LICENSE)
[![Go Report Card](https://img.shields.io/github/go-mod/go-version/bensadeh/circumflex?color=ffe585&labelColor=292D3E)](https://github.com/bensadeh/circumflex/blob/master/go.mod)
</div>

`circumflex` is a command line tool for browsing Hacker&nbsp;News in your terminal.

<p align="center">
  <img src="screenshots/mainview.png" width="700" alt="^"/>
</p>

## Getting started

Install `circumflex` with Homebrew:

```console
# Install
brew install bensadeh/circumflex/circumflex

# Run
clx
```

Press <kbd>i</kbd> to show available keymaps and settings.

## Features

### Overview

`circumflex` is a text-based user interface (TUI) application that lets you browse Hacker News in your terminal. It can
list submissions by category and show the comment section. 

On the Main View, `circumflex` supports Vim keybindings for navigating between submissions, including <kbd>h</kbd>, 
<kbd>j</kbd>, <kbd>k</kbd> and <kbd>l</kbd>, digits and <kbd>g</kbd>/<kbd>G</kbd>.

`circumflex` is made to look well across different configurations by adhering to your terminal's native color scheme. 


<p align="center">
  <img src="screenshots/comments.png" width="700" alt="^"/>
</p>



### Comment section

Comments are pretty-printed and piped to the pager `less`. To present a nice and readable comment section,
`circumflex` features:

* Text in **bold**, _italics_, [hyperlinks](https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda) and
  `code` where available
* Labels for Original Posters (`OP`), Parent Posters (`PP`) and moderators (`mod`)
* Ability to jump between top-level comments by searching for `::`

### Syntax highlighting
#### Quotes

In the comment section in Hacker News, quotes are separated by a paragraph and open with a `>`. However, besides these 
distinctions, quotes are not further stylized or formatted:
```
> This is a quote 

This is not a quote
```

In `circumflex`, the `>` symbol is omitted and quotes are instead italicised and dimmed.

<p align="center">
  <img src="screenshots/quotes.png" width="700" alt="^"/>
</p>

#### Headlines
Certain submissions on the front page are also in their own 'category', such as `Show HN` and `Ask HN`. `circumflex` 
colorizes these parts of the headlines to make it easier to distinguish them from each other. On by default, but can be 
disabled.

<p align="center">
  <img src="screenshots/showtell.png" width="250" alt="^"/>
</p>

#### YC-funded startups
[Twice a year](https://www.ycombinator.com/companies/), Y Combinator funds start-ups through its accelerator program. 
`circumflex` highlights these startups to signalize their affiliation with YC. On by default, but can be disabled.

<p align="center">
  <img src="screenshots/yc.png" width="350" alt="^"/>
</p>

#### References
References are formatted as numbers inside brackets: `[1]`, `[2]`, `[…]`. `circumflex` highlights the numbers 
inside the brackets for easier cross-referencing.

<p align="center">
  <img src="screenshots/linkHighlights.png" width="700" alt="^"/>
</p>

## Favorites
You can save submissions you'd like to revisit later by adding them to Favorites. Press <kbd>f</kbd> to add the 
currently highlighted submission to your Favorites list. Favorited items can be removed with <kbd>x</kbd>.

Favorites are stored in `favorites.json` in `~/.config/circumflex`.

## Settings

### How to configure
#### Overview
The available options and their current values are displayed in the Settings View. Overridden values are marked with
`*`. To enter Settings View, press <kbd>i</kbd> on the Main View and then <kbd>Tab</kbd> to change the category. Note 
that the settings cannot be changed from within the Settings View.

<p align="center">
  <img src="screenshots/settings.png" width="700" alt="^"/>
</p>

To change the settings, you can either:

1. edit `config.env` in `~/.config/circumflex`, or
2. set environment variables in your shell

#### Changing settings through `config.env`
`circumflex` can create a `config.env` in `~/.config/circumflex` by pressing <kbd>t</kbd> in the Settings View. This
file will contain the available options, their descriptions and default values.

#### Changing settings with environment variables
Depending on your preference, it might be more convenient for you to configure `circumflex` by setting 
[environment variables](https://unix.stackexchange.com/questions/117467/how-to-permanently-set-environmental-variables).
Below are a couple of examples on how to set the variables in different shells. Run the commands directly from your 
shell to set the value for the current session. Put the commands in somewhere in your 
[dotfiles](https://dotfiles.github.io/) to make the settings permanent.

Bash and zsh:
```bash
export CLX_COMMENT_WIDTH=65
```

Fish:
```fish
set -x CLX_COMMENT_WIDTH "65"
```

### Available options

The following table shows the ways `circumflex` can be configured:

| Key                         | Default Value | Description |
| :-------------------------: | :-------: |---|
| `CLX_COMMENT_WIDTH`         | `70` | Sets the maximum number of characters on each line for comments, replies and descriptions in settings. Set to 0 to use the whole screen.       |
| `CLX_INDENT_SIZE`           | `4` | The number of whitespaces prepended to each reply, not including the color bar.        |
| `CLX_HIGHLIGHT_HEADLINES`   | `2` | Highlights YC-funded startups and text containing `Show HN`, `Ask HN`, `Tell HN` and `Launch HN`. Can be set to 0 (No highlighting), 1 (inverse highlighting) or 2 (colored highlighting).        |
| `CLX_RELATIVE_NUMBERING`    | `false` | Shows each line with a number relative to the currently selected element. Similar to Vim's hybrid line number mode.        |
| `CLX_HIDE_YC_JOBS`          | `true` | Hides `X is hiring` posts from YC-funded startups. Does not affect the monthly `Who is Hiring?` posts.        |
| `CLX_PRESERVE_RIGHT_MARGIN` | `false` | Shortens replies so that the total length, including indentation, is the same as the comment width. Best used when Indent Size is small to avoid deep replies being too short.   |

## Under the hood

`circumflex` uses:

* [tcell](https://github.com/gdamore/tcell) and [cview](https://gitlab.com/tslocum/cview) for the TUI
* [viper](https://github.com/spf13/viper) for reading and setting configurations
* [cheeaun's unofficial Hacker News API](https://github.com/cheeaun/node-hnapi) for providing submissions and comments
* [`less`](http://greenwoodsoftware.com/less/) for viewing comments
* [go-term-text](https://github.com/MichaelMure/go-term-text) for text formatting

Screenshots use:

* [iTerm2](https://iterm2.com/) for the terminal
* [Palenight Theme](https://github.com/JonathanSpeek/palenight-iterm2) for the color scheme
* [JetBrains Mono](https://github.com/JetBrains/JetBrainsMono) for the font
