package music

import (
	"github.com/ItsClairton/Anny/core"
	commands "github.com/ItsClairton/Anny/music/commands"
	events "github.com/ItsClairton/Anny/music/events"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Module = &core.Module{
	Name: "Música", Emote: emojis.Yeah,
	Commands: []*core.Command{&commands.PlayCommand, &commands.SkipCommand, &commands.StopCommand, &commands.PauseCommand, &commands.ResumeCommand, &commands.SeekCommand, &commands.NowplayingCommand, &commands.ShuffleCommand, &commands.QueueCommand},
	Events:   []*core.Event{&events.VServerUpdateEvent, &events.VStateUpdateEvent},
}
