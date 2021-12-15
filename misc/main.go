package misc

import (
	"github.com/ItsClairton/Anny/core"
	commands "github.com/ItsClairton/Anny/misc/commands"
	events "github.com/ItsClairton/Anny/misc/events"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Module = &core.Module{
	Name: "Miscelânea", Emote: emojis.Peer,
	Commands: []*core.Command{&commands.PingCommand, &commands.HelpCommand},
	Events:   []*core.Event{&events.ReadyEvent},
}
