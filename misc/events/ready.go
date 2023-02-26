package misc

import (
	"github.com/ItsClairton/Anny/core"
	"github.com/diamondburned/arikawa/v3/gateway"
)

var ReadyEvent = core.Event{
	Handler: func(_ *gateway.ReadyEvent) {
		// core.State.UpdateStatus(gateway.UpdateStatusData{Activities: []discord.Activity{{Name: os.Getenv("DISCORD_STATUS"), Type: discord.ListeningActivity}}})
	},
}
