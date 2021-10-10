package music

import (
	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var Category = &discord.Category{
	Name:         "Música",
	Emote:        emojis.PingPong,
	Interactions: []*discord.Interaction{&SkipCommand, &PauseCommand, &UnpauseCommand, &NowplayingCommand, &ShuffleCommand, &StopCommand},
}
