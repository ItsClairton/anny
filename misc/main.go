package misc

import (
	"github.com/ItsClairton/Anny/core"
	commands "github.com/ItsClairton/Anny/misc/commands"
	events "github.com/ItsClairton/Anny/misc/events"
)

var Module = &core.Module{
	Commands: []*core.Command{commands.PingCommand},
	Events:   []*core.Event{events.ReadyEvent},
}
