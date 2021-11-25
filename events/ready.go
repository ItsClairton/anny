package events

import (
	"os"

	"github.com/ItsClairton/Anny/core"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func OnReady(_ *gateway.ReadyEvent) {
	core.Session.UpdateStatus(gateway.UpdateStatusData{
		Activities: []discord.Activity{{
			Name: os.Getenv("DISCORD_STATUS"),
			Type: discord.ListeningActivity,
		}},
	})
}
