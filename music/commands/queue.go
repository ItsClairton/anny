package music

import (
	"time"

	"github.com/ItsClairton/Anny/core"
	music "github.com/ItsClairton/Anny/music/audio"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
)

var QueueCommand = core.Command{
	Name:        "fila",
	Description: "Mostra as músicas adicionadas na fila",
	Handler: func(ctx *core.CommandContext) {
		player := music.GetPlayer(ctx.GuildID)

		if player == nil || len(player.Queue) == 0 {
			ctx.Ephemeral().Reply(emojis.Cry, "Não há nada na fila no momento.")
			return
		}

		var finalText string
		var count int

		for _, track := range player.Queue {
			if count > 9 {
				break
			}

			count++
			finalText += utils.Fmt("%s - [%s](%s)\n", emojis.GetNumberAsEmoji(count), track.Title, track.URL)
		}

		ctx.Reply(utils.NewEmbed().
			Color(0xA652BB).
			Description(finalText).
			Footer(utils.Fmt("Mostrando %d de %d músicas", count, len(player.Queue)), ctx.User.AvatarURL()).
			Timestamp(time.Now()))
	},
}
