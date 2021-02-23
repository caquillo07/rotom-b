# Rotom-B Pokémon Discord bot 

Rotom-B is a Discord Bot with lots of easily accessible Pokémon information, tailored to make your Max Raid battles easier. Rotom-B keeps up to date Pokémon Sword & Shield data including: dens, catch rate, pokéballs, pokédex, types and natures data, with more to be added soon.

## Top.gg
[![Discord Bots](https://top.gg/api/widget/726478988276531212.svg)](https://top.gg/bot/726478988276531212)

[Vote for us on Top.gg!](https://top.gg/bot/726478988276531212)


## How to get started

[Invite Rotom-B](https://discord.com/oauth2/authorize?client_id=703830812772139138&scope=bot&permissions=281600) to your Discord Server to get started!

Rotom-B will be added with simple message reading permissions and will start listening to commands with the default `$` prefix.

## Why Rotom-B? 
Rotom-B is a rewrite of the popular [Alcremie-B](https://github.com/RaphGG/den-bot) den bot, with some extra features.

We loved Alcremie-B! Unfortunately the maintainer was unable to keep maintaining the bot, so we decided it would be a fun project to rewrite it, and add some cool features along the way.
This is not meant to be competition, so please all credits goes to the original creator [RaphGG](https://github.com/RaphGG) for the idea. 

## Features
- Complete Max Raid Dens information, including up to date Isle of Armor DLC. 
- Easily accessible Pokémon data, like stats, abilities, den locations, pokédex entries, etc. 
- Catch rate & pokéball stats calculator based on Bulbapedia.
- Sprites for all Pokémon, including forms and shiny versions. 

## Commands
-  `<>` indicate required fields.
-  `[]` indicate optional fields.
- Use `*` next to Pokémon name for shiny sprites.
- Catch Rates calculation are under Raid Specific Conditions: Levels 30-70, 1 HP, and no status modifiers.

Command | Arguments | Description
--- | --- | ---
`$ball` | `<ball_name>` | Shows a summary of a Poké-Ball’s statistics
`$catch` | `<pokemon> [form] [ball_name]` | Summary of catch rates for a given Pokémon and Ball combination.
`$credits` | | Credits to all who helped in the creation of the bot.
`$den` | `<den_number/pokemon_name>` | Shows a list of Pokémon that belong to a den including their HAs.
`$help` | | Displays a list of commands you have access to use.
`$invite`| | Get an Invite Link to invite Rotom-B to another server!
`$nature`| `<nature>` | Shows ithe Pokémon Sprite in appropriate form
`$pokedex` | `<pokemon>`| Shows Pokédex info on every Pokémon.
`$sprite` |  `<pokemon>` |  Shows the Pokémon Sprite. Include * in the end for the shiny sprite.
`$type` | `<type>` | Shows info regarding Pokémon Types.
`$version` |  | Check which version of Rotom-B is running.
`$settings` | `[setting] <new_value>` | Allows administrators to set server specific configuration.

## Upcoming Features/Todos
- [X] ~~Persistency, including a database for all the data~~
- [X] ~~Custom settings~~
- [X] ~~Updated Pokémon information with all new Sword and Shield DLC Pokémon and Zarude~~
- [ ] Automatic messages on bot status and new updates
- [ ] Friend codes and IGN storage support
- [ ] Creating website
- [ ] ...And we're accepting ideas/PRs! :) 

## Observations and Know Issues
Please note we don't have an animated sprite for Gigantamax Inteleon and its Shiny version. This is the only Pokémon we're using a still image for a sprite. Pull requests adding it are much appreciated!

All images are hosted in [this repository](https://github.com/caquillo07/rotom-b-data) 

## Screenshots
What the bot looks like on Discord!

![Pokémon Sprites](https://raw.githubusercontent.com/hypermilla/caquillo07.github.io/master/rotomb_screenshots/rotomB_pkmn_sprite.png)

![Catch Rates](https://raw.githubusercontent.com/hypermilla/caquillo07.github.io/master/rotomb_screenshots/rotomB_catchrates.png)

![Den Pokémon Information](https://raw.githubusercontent.com/hypermilla/caquillo07.github.io/master/rotomb_screenshots/rotomB_pkmn_den.png)

![Den Information from number](https://raw.githubusercontent.com/hypermilla/caquillo07.github.io/master/rotomb_screenshots/rotomB_den_number.png)


## Development/Deployment

### Prerequisites

* [Go](https://golang.org/) - Go Programming language
* [PostgreSQL](https://www.postgresql.org/download/) - A working database instance ([Docker recommended](https://hub.docker.com/_/postgres))
* [Air - Live Reload](https://github.com/cosmtrek/air) - (optional) if running live reload while development.
* [Docker](https://docs.docker.com/get-docker/) - (optional) if running the database on Docker.

First we need a configuration file, copy the example-config.yaml and update values as needed
```shell script
cp example-config.yaml config.yaml
```

The Simplest way to get everything running is using `docker-compose`
```shell script
docker-compose up
```

If running individually, first run a postgres instance:
```shell script
# Run the postgres instance
docker run --name postgres \ 
    -e POSTGRES_PASSWORD=root \ 
    -e POSTGRES_USER=postgres \ 
    -e POSTGRES_DB=postgres \ 
    -p 5432:5432 \ 
    -d postgres

# Run migrations
make migrate-dev
```

For simple development, you can just do the following:
```shell script
make dev
```

For auto-reload on changes, you can run this instead (requires [air installed](https://github.com/cosmtrek/air#installation))
```shell script
make dev-reload
```

For a production build, you must run the make command for specific OS needs. At the moment we have support for AMD64 Linux/OSX/Windows (windows is untested).
```shell script
# Creating production build for Linux
make linux

# Run the bot
./den-bot-linux-amd64 bot
```

For a production build on any other architecture not in the Makefile (PRs welcomed!)
```shell script
# Setup the ENV variables needed
OS=<desired_os>
GOARCH=<desired_arch>
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
VERSION=$(shell cat .version)
METRICS_IMPORT_PATH=github.com/caquillo07/rotom-bot/metrics

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X ${METRICS_IMPORT_PATH}.Version=${VERSION} -X ${METRICS_IMPORT_PATH}.Commit=${COMMIT} -X ${METRICS_IMPORT_PATH}.Branch=${BRANCH}"

# Build the bot
GOOS=${OS} GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-${OS}-${GOARCH} .
```

# Built With
* [Go](https://golang.org/) - Programming language
* [cobra](https://github.com/spf13/cobra) - CLI
* [viper](https://github.com/spf13/viper) - Configuration
* [zap](https://github.com/uber-go/zap) - Logging internals
* [DiscordGo](https://github.com/bwmarrin/discordgo) - Discord Client
* [GORM](https://github.com/go-gorm/gorm) - The fantastic ORM library for Go
* [go-cache](https://github.com/patrickmn/go-cache) - In memory cache for Go - Similar to Memcached
