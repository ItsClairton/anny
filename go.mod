module github.com/ItsClairton/Anny

go 1.17

require (
	github.com/Pauloo27/searchtube v0.0.0-20210906001334-44c3e43c257a
	github.com/bwmarrin/discordgo v0.23.3-0.20211010150959-f0b7e81468f7
	github.com/jonas747/ogg v0.0.0-20161220051205-b4f6f4cf3757
	github.com/kkdai/youtube/v2 v2.7.4
	github.com/buger/jsonparser v1.1.1
	github.com/joho/godotenv v1.4.0
)

require (
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
)

replace github.com/kkdai/youtube/v2 => github.com/xesnault/youtube/v2 v2.7.5-0.20211016224312-4057b33ef4cf
