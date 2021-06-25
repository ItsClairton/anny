# Anny
Simple bot for Dscord witten in Golang
## Commands
| Name      | Description                                                                                                           | Category      |
|-----------|-----------------------------------------------------------------------------------------------------------------------|---------------|
|`ping`     | Respond you with bot latency                                                                                          | Miscellaneous |
|`help`     | Respond with a list of commands                                                                                       | Miscellaneous |
|`scene`    | Trace episode, name of anime and time of a matching scene from a screenshot using [trace.moe](https://trace.moe/about)| Utilities     |
|`anime`    | Shows basic information of an anime (Using AniList, MAL and Google Translate)                                         | Utilities     |
|`manga`    | Shows basic information of an manga (Using AniList, MAL and Google Translate)                                         | Utilities     |
|`translate`| Translate text to another language (Using Google Translate)                                                           | Utilities     |
|`cat`      | Generate random cat images (Using [TheCatAPI](https://thecatapi.com/) and [NekosLife](https://nekos.life/))           | Image         |
|`neko`     | Generate random neko images (Using [NekosLife](https://nekos.life/))                                                  | Image         |

## WARNING
- This bot was created while I was learning about Golang, so it can have a lot of bugs, Sorry for my English.

## Selfhost
- Create a .env file and copy the content of the .env.example file into it.
- To build use `go build .`
- To run use `go run .`, or `./Anny` if you already builded