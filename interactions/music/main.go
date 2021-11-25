package music

import (
	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &core.Category{
	Name:         "Música",
	Emote:        emojis.PingPong,
	Interactions: []*core.Interaction{&PlayCommand, &SkipCommand, &PauseCommand, &ResumeCommand, &ShuffleCommand, &NowplayingCommand, &StopCommand, &SeekCommand},
}
