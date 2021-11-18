package events

import (
	"os"

	"github.com/ItsClairton/Anny/base"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func OnReady(_ *gateway.ReadyEvent) {
	base.Session.UpdateStatus(gateway.UpdateStatusData{
		Activities: []discord.Activity{{
			Name: os.Getenv("DISCORD_STATUS"),
			Type: discord.ListeningActivity,
		}},
	})
}
