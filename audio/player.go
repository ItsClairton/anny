package audio

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/voice"
)

var (
	StoppedState   = 0
	LoadingState   = 1
	PausedState    = 2
	PlayingState   = 3
	DestroyedState = 4

	players = map[discord.GuildID]*Player{}
)

type Player struct {
	*sync.Mutex

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
	player := &Player{Mutex: &sync.Mutex{}, State: StoppedState, GuildID: GuildID, TextID: TextID, VoiceID: VoiceID}
	players[player.GuildID] = player

	go func() {
		s, _ := voice.NewSession(base.Session)

		err := s.JoinChannel(GuildID, VoiceID, false, true)
		if err != nil {
			player.Kill(true, emojis.Cry, "Um erro ocorreu ao tentar se conectar ao canal de voz: `%v`", err)
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
	if p.State != StoppedState {
		return
	}

	if len(p.Queue) == 0 || p.Connection == nil {
		go p.Kill(false)
		return
	}

	p.State = LoadingState
	if p.Timer != nil {
		p.Timer.Stop()
		p.Timer = nil
	}

	current := p.Queue[0]
	p.Queue = p.Queue[1:]

	if !current.IsLoaded() {
		if song, err := current.Load(); err != nil {
			p.State = StoppedState
			go p.Play()

			base.SendMessage(p.TextID, emojis.Cry, "Um erro ocorreu ao carregar a música **%s**: `%v`", current.Title, err)
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
			base.SendMessage(p.TextID, emojis.Cry, "Um erro ocorreu enquanto tocava a música **%s**: `%v`", current.Title, err)
		}

		p.Current, p.State = nil, StoppedState
		p.Play()
	}()

	embed := base.NewEmbed().
		SetDescription("%s Tocando agora [%s](%s)", emojis.Yeah, current.Title, current.URL).
		SetImage(current.Thumbnail).
		SetColor(0xA652BB).
		AddField("Autor", current.Author, true).
		AddField("Duração", utils.Is(current.IsLive, "--:--", utils.FormatTime(current.Duration)), true).
		AddField("Provedor", current.Provider(), true).
		SetFooter(utils.Fmt("Pedido por %s#%s", current.Requester.Username, current.Requester.Discriminator), current.Requester.AvatarURL()).
		SetTimestamp(current.Time)

	base.Session.SendMessage(p.TextID, "", embed.Build())
}

func (p *Player) Kill(force bool, args ...interface{}) {
	p.Lock()

	removePlayer := func() {
		p.Lock()

		if p.Timer != nil {
			p.Timer.Stop()
			p.Timer = nil
		}

		if force || p.State == StoppedState {
			println(len(args))
			if len(args) >= 2 {
				base.SendMessage(p.TextID, args[0].(string), args[1].(string), args[2:]...)
			}

			p.Queue = []*RequestedSong{}

			if p.Current != nil {
				p.Current.Stop()
			}

			if p.Connection != nil {
				p.Connection.Leave()
			}

			p.State = DestroyedState
			delete(players, p.GuildID)
		}

		p.Unlock()
	}

	if force {
		p.Unlock()
		removePlayer()
		return
	}

	if p.Timer == nil && p.State == StoppedState {
		p.Timer = time.AfterFunc(3*time.Minute, removePlayer)
	}

	p.Unlock()
}
