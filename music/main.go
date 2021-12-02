package music

import (
	"github.com/ItsClairton/Anny/core"
	commands "github.com/ItsClairton/Anny/music/commands"
)

var Module = &core.Module{
	Commands: []*core.Command{&commands.PlayCommand, &commands.SkipCommand, &commands.StopCommand, &commands.PauseCommand, &commands.ResumeCommand, &commands.SeekCommand, &commands.NowplayingCommand},
	Events:   []*core.Event{},
}
