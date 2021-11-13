package audio

import (
	"io"
	"math/rand"
	"time"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/voice"
)

var (
	StoppedState = 0
	PausedState  = 1
	PlayingState = 2

	players = map[discord.GuildID]*Player{}
)

type Player struct {
	Timer      *time.Timer
	Connection *voice.Session

	State   int
	Current *CurrentSong
	Queue   []*RequestedSong

	GuildID         discord.GuildID
	VoiceID, TextID discord.ChannelID
}

type RequestedSong struct {
	*Song
	Requester *discord.User
	Time      time.Time
}

type CurrentSong struct {
	*RequestedSong
	*StreamingSession
}

func NewPlayer(GuildID discord.GuildID, TextID, VoiceID discord.ChannelID) *Player {
	player := &Player{State: StoppedState, GuildID: GuildID, TextID: TextID, VoiceID: VoiceID}
	players[player.GuildID] = player

	go func() {
		s, _ := voice.NewSession(base.Session)

		err := s.JoinChannel(GuildID, VoiceID, false, true)
		if err != nil {
			base.SendMessage(player.TextID, emojis.MikuCry, "Um erro ocorreu ao tentar se conectar ao canal de voz: `%v`", err)
			player.Kill(true)
		}

		player.Connection = s
		player.Play()
	}()

	return player
}

func GetPlayer(id discord.GuildID) *Player {
	return players[id]
}

func (p *Player) AddSong(requester *discord.User, shuffle bool, tracks ...*Song) {
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
			base.SendMessage(p.TextID, emojis.MikuCry, "Um erro ocorreu ao carregar a música **%s**: `%v`", current.Title, err)
			go p.Play()
			return
		} else {
			current.Song = song
		}
	}

	done := make(chan error)
	p.Current, p.State = &CurrentSong{current, StreamURL(current.StreamingURL, p.Connection, done)}, PlayingState
	go func() {
		err := <-done

		if err != io.EOF {
			base.SendMessage(p.TextID, emojis.MikuCry, "Um erro ocorreu enquanto tocava a música **%s**: `%v`", current.Title, err)
		}

		p.Current, p.State = nil, StoppedState
		p.Play()
	}()

	embed := base.NewEmbed().
		SetDescription("%s Tocando agora [%s](%s)", emojis.ZeroYeah, current.Title, current.URL).
		SetImage(current.Thumbnail).
		SetColor(0xA652BB).
		AddField("Autor", current.Author, true).
		AddField("Duração", utils.Is(current.IsLive, "--:--", utils.FormatTime(current.Duration)), true).
		AddField("Provedor", current.Provider(), true).
		SetFooter(utils.Fmt("Pedido por %s#%s", current.Requester.Username, current.Requester.Discriminator), current.Requester.AvatarURL()).
		SetTimestamp(current.Time)

	base.Session.SendMessage(p.TextID, "", embed.Build())
}

func (p *Player) Kill(force bool) {
	removePlayer := func() {
		p.Queue = []*RequestedSong{}

		if p.Current != nil {
			p.Current.Stop()
		}

		if p.Connection != nil {
			p.Connection.Leave()
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
