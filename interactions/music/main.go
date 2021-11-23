package music

import (
	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &base.Category{
	Name:         "Música",
	Emote:        emojis.PingPong,
	Interactions: []*base.Interaction{&PlayCommand, &SkipCommand, &PauseCommand, &ResumeCommand, &ShuffleCommand, &NowplayingCommand, &StopCommand, &SeekCommand},
}
