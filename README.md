# Rotom-B Pokémon Discord bot 

Rotom-B is a Discord Bot with lots of easily accessible Pokémon information, tailored to make your Max Raid battles easier. Rotom-B keeps up to date Pokémon Sword & Shield data including: dens, catch rate, pokéballs, pokédex, types and natures data, with more to be added soon.


## How to get started

[Invite Rotom-B](https://discord.com/oauth2/authorize?client_id=703830812772139138&scope=bot&permissions=281600) to your Discord Server to get started!

Rotom-B will be added with simple message reading permissions and will start listening to commands with the default `$` prefix.

## Why Rotom-B? 
Rotom-B is a rewrite of the popular Alcremie-B den bot, with some extra features.  

## Features
- Complete Max Raid Dens information, including up to date Isle of Armor DLC. 
- Easily accessible Pokémon data, like stats, abilities, den locations, pokédex entries, etc. 
- Catch rate & pokéball stats calculator based on Bulbapedia.
- Sprites for all pokémon, including forms and shiny versions. 

## Commands
-  `<>` indicate required fields.
-  `[]` indicate optional fields.
- Use `*` for shiny sprites.
- Catch Rates are calculated under Raid Specific Conditions: Levels 30-70, 1 HP, and no status modifiers.

Command | Arguments | Description
--- | --- | ---
`$ball` | `<ball_name>` | Shows a summary of a Poké-Ball’s statistics
`$catch` | `<pokemon> [form] [ball_name]` | Summary of catch rates for a given Pokémon and Ball combination.
`$credits` | | Credits to all who helped in the creation of the bot.
`$den` | `<den_number|pokemon_name>` | Shows a list of Pokémon that belong to a den including their HAs.
`$invite`| | Get an Invite Link to invite Rotom-B to another server!
`$nature`| `<nature>` | Shows ithe Pokémon Sprite in appropriate form
`$pokedex` | `<pokemon>`| Shows Pokédex info on every Pokémon.
`$sprite` |  `<pokemon>` |  Shows the Pokémon Sprite. Include * in the end for the shiny sprite.
`$type` |  `<type>` | Shows info regarding Pokémon Types.
` $version`|  | Check which version of Rotom-B is running.

## Upcoming Features
- Persistency, including a database for all the data
- Custom settings
- Updated Sword and Shield sprites for all Pokémon
- ...And we're accepting ideas! :) 

## Observations
Please note we don't have an animated sprite for Gigantamax Inteleon and its Shiny version. This is the only Pokémon we're using a still image for a sprite. Pull requests adding it are much appreciated!

All images are hosted in [this repository](https://github.com/caquillo07/rotom-b-data) 

## Screenshots
What the bot looks like on Discord!

![Pokémon Sprites](https://raw.githubusercontent.com/hypermilla/caquillo07.github.io/master/rotomb_screenshots/rotomB_pkmn_sprite.png)

![Catch Rates](https://raw.githubusercontent.com/hypermilla/caquillo07.github.io/master/rotomb_screenshots/rotomB_catchrates.png)

![Den Pokémon Information](https://raw.githubusercontent.com/hypermilla/caquillo07.github.io/master/rotomb_screenshots/rotomB_pkmn_den.png)

![Den Information from number](https://raw.githubusercontent.com/hypermilla/caquillo07.github.io/master/rotomb_screenshots/rotomB_den_number.png)






