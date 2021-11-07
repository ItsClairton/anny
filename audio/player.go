package audio

import (
	"io"
	"math/rand"
	"time"

	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	StoppedState = 0
	PausedState  = 1
	PlayingState = 2

	players = map[string]*Player{}
)

type Player struct {
	Timer      *time.Timer
	Connection *discordgo.VoiceConnection

	State   int
	Current *CurrentSong
	Queue   []*RequestedSong

	GuildID, TextID, VoiceID string
}

type RequestedSong struct {
	*Song
	Requester *discordgo.User
	Time      time.Time
}

type CurrentSong struct {
	*RequestedSong
	*StreamingSession
}

func NewPlayer(GuildID, TextID, VoiceID string) *Player {
	player := &Player{State: StoppedState, GuildID: GuildID, TextID: TextID, VoiceID: VoiceID}
	players[player.GuildID] = player

	go func() {
		connection, err := discord.Session.ChannelVoiceJoin(player.GuildID, player.VoiceID, false, true)
		if err != nil {
			discord.SendMessage(player.TextID, emojis.MikuCry, "Um erro ocorreu ao fazer conexão com o Canal de voz: `%v`", err)
			logger.Warn(utils.Fmt("Um erro ocorreu ao fazer conexão com o canal de voz %s, da Guilda %s.", player.VoiceID, player.GuildID), err)
			player.Kill(true)
		} else {
			player.Connection = connection
			player.Play()
		}
	}()

	return player
}

func GetPlayer(id string) *Player {
	return players[id]
}

func (p *Player) AddSong(requester *discordgo.User, shuffle bool, tracks ...*Song) {
	for _, track := range tracks {
		p.Queue = append(p.Queue, &RequestedSong{track, requester, time.Now()})
	}

	if shuffle && len(tracks) > 1 {
		p.Shuffle()
	}

	go p.Play()
}

func (p *Player) Skip() {
	p.Current.Stop()
}

func (p *Player) Pause() {
	p.Current.Pause(true)
	p.State = PausedState
}

func (p *Player) Resume() {
	p.Current.Pause(false)
	p.State = PlayingState
}

func (p *Player) Shuffle() {
	rand.Shuffle(len(p.Queue), func(old, new int) {
		p.Queue[old], p.Queue[new] = p.Queue[new], p.Queue[old]
	})
}

func (p *Player) Play() {
	if p.State != StoppedState || players[p.GuildID] == nil {
		return
	}

	if len(p.Queue) == 0 || p.Connection == nil {
		go p.Kill(false)
		return
	}

	if p.Timer != nil {
		p.Timer.Stop()
		p.Timer = nil
	}

	current := p.Queue[0]
	p.Queue = p.Queue[1:]

	if !current.IsLoaded() {
		if song, err := current.Load(); err != nil {
			discord.SendMessage(p.TextID, emojis.MikuCry, "Um erro ocorreu ao carregar a música **%s**: `%v`", current.Title, err)
			go p.Play()
			return
		} else {
			current.Song = song
		}
	}

	done := make(chan error)
	p.Current, p.State = &CurrentSong{current, StreamFromPath(current.StreamingURL, p.Connection, done)}, PlayingState
	go func() {
		err := <-done

		if err == ErrVoiceTimeout {
			discord.SendMessage(p.TextID, emojis.MikuCry, "Tempo de conexão com o canal de voz esgotado.")
			p.Current, p.State = nil, StoppedState
			p.Kill(true)
			return
		}

		if err != io.EOF {
			discord.SendMessage(p.TextID, emojis.MikuCry, "Um erro ocorreu enquanto tocava a música **%s**: `%v`", current.Title, err)
		}

		p.Current, p.State = nil, StoppedState
		p.Play()
	}()

	embed := discord.NewEmbed().
		SetDescription("%s Tocando agora [%s](%s)", emojis.ZeroYeah, current.Title, current.URL).
		SetImage(current.Thumbnail).
		SetColor(0xA652BB).
		AddField("Autor", current.Author, true).
		AddField("Duração", utils.Is(current.IsLive, "--:--", utils.FormatTime(current.Duration)), true).
		AddField("Provedor", current.Provider(), true).
		SetFooter(utils.Fmt("Pedido por %s#%s", current.Requester.Username, current.Requester.Discriminator), current.Requester.AvatarURL("")).
		SetTimestamp(current.Time.Format(time.RFC3339))

	discord.Session.ChannelMessageSendEmbed(p.TextID, embed.Build())
}

func (p *Player) Kill(force bool) {
	removePlayer := func() {
		p.Queue = []*RequestedSong{}

		if p.Current != nil {
			p.Current.Stop()
		}

		if p.Connection != nil {
			p.Connection.Disconnect()
		}

		delete(players, p.GuildID)
	}

	if force {
		removePlayer()
		return
	}

	if p.Timer == nil && p.State == StoppedState && len(p.Queue) == 0 {
		p.Timer = time.AfterFunc(5*time.Minute, removePlayer)
	}
}
