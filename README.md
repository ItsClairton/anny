# Anny
Simple bot for Dscord witten in Golang
## Commands
| Name  | Description                                                                                                           | Category      |
|-------|-----------------------------------------------------------------------------------------------------------------------|---------------|
|`ping` | Respond you with bot latency                                                                                          | Miscellaneous |
|`scene`| Trace episode, name of anime and time of a matching scene from a screenshot using [trace.moe](https://trace.moe/about)| Anime         |
|`anime`| Shows basic information of an anime (Uses AniList, MAL and Google Translate)                                          | Anime         |
|`manga`| Shows basic information of an manga (Uses AniList, MAL and Google Translate)                                          | Anime         |
|`cat`  | Generate random cat images (Using [TheCatAPI](https://thecatapi.com/) and [NekosLife](https://nekos.life/))           | Image         |
|`neko` | Generate random neko images (Using [NekosLife](https://nekos.life/))                                                  | Image         |

## WARNING
- This bot was created while I was learning about Golang, so it can have a lot of bugs, Sorry for my English.

## Selfhost
- Create a .env file and copy the content of the .env.example file into it.
- To build use `go build .`
- To run use `go run .`, or `./Anny` if you already builded