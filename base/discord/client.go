package discord

import "github.com/bwmarrin/discordgo"

var (
	Session *discordgo.Session
)

func Init(token string) {
	Session, _ = discordgo.New("Bot " + token)

	Session.Identify.Intents = discordgo.IntentsAll
}

func Connect() error {
	return Session.Open()
}

func Disconnect() {
	Session.Close()
}
