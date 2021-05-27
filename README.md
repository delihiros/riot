# Riot

"riot" is both a library and a CLI for Riot Games' API.

Please visit [Riot Developer Portal](https://developer.riotgames.com/apis) for more information about its API.


## Current status:
- [x] Valorant API Support
  - [x] CLI
- [ ] LoL API Support
- [ ] LoR API Support

## Usage:

```sh
Riot is a command line interface for RiotGames API

Usage:
  riot [command]

Available Commands:
  help        Help about any command
  valorant    subcommand for valorant

Flags:
      --api-key string   riot api key
      --config string    config file (default is $HOME/.riot.yaml)
  -h, --help             help for riot
      --region string    region (default "ap")

Use "riot [command] --help" for more information about a command.
```