package base

import (
	"context"

	"github.com/ItsClairton/Anny/utils"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var (
	Session *state.State
)

func New(token string) error {
	Session = state.NewWithIntents(utils.Fmt("Bot %s", token), gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentGuildVoiceStates)

	return Session.Open(context.Background())
}

func Me() *discord.User {
	me, _ := Session.Me()

	return me
}

func SendMessage(id discord.ChannelID, emoji, text string, args ...interface{}) {
	Session.SendMessage(id, utils.Fmt("%s | %s", emoji, utils.Fmt(text, args...)))
}

func Disconnect() {
	Session.Close()
}
