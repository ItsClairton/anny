package music

var players = map[string]*Player{}

func GetPlayer(guildId string) *Player {

	result, exist := players[guildId]

	if !exist {
		return nil
	}

	return result
}

func AddPlayer(player *Player) *Player {
	players[player.GuildID] = player
	return player
}
