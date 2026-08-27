# Riot

`riot` is a Go library and command-line interface (CLI) for the Riot Games API.

See the [Riot Developer Portal](https://developer.riotgames.com/apis) for API details.

## Supported APIs

The following APIs are supported:

- [x] VALORANT PC API
  - [x] Content, Match, Ranked, and Status clients
  - [x] CLI, including recent matches
- [x] VALORANT Console Match and Ranked clients
- [x] Account API, including Riot Sign On (RSO) `/riot/account/v1/accounts/me`
- [ ] LoL API
- [x] LoR Deck, Inventory, Match, Ranked, and Status clients

## Install the CLI

Go 1.26.7 or newer is required.

```sh
go install github.com/delihiros/riot@latest
```

## Configure and use the CLI

Set the API key with an environment variable (recommended):

```sh
export RIOT_API_KEY=RGAPI-your_api_key_here
riot valorant status get-status --region ap
```

Alternatively, put `riot.api_key` in `$HOME/.riot.yaml`, pass another file with `--config`, or override it with `--api-key`.

Run `riot --help` to list the available commands and global flags:

```sh
Riot is a command line interface for RiotGames API

Usage:
  riot [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  valorant    subcommand for valorant

Flags:
      --api-key string   riot api key override
      --config string    config file (default is $HOME/.riot.yaml)
  -h, --help             help for riot
      --region string    region (default "ap")

Use "riot [command] --help" for more information about a command.
```

`valorant content get-content` returns unfiltered content by default. Pass `--locale`, such as `--locale ja-JP`, to request one locale.
