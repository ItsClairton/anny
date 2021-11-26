package audio

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
)

const (
	StoppedState = iota
	LoadingState
	PausedState
	PlayingState
	DestroyedState
)

var players = map[discord.GuildID]*Player{}

type Player struct {
	*sync.Mutex

	Timer   *time.Timer
	Session *VoicySession

	State   int
	Current *RequestedSong
	Queue   []*RequestedSong

	GuildID         discord.GuildID
	VoiceID, TextID discord.ChannelID
}

type RequestedSong struct {
	*Song
	Requester *discord.User
	Time      time.Time
}

func NewPlayer(GuildID discord.GuildID, TextID, VoiceID discord.ChannelID) *Player {
	player := &Player{Mutex: &sync.Mutex{}, State: StoppedState, GuildID: GuildID, TextID: TextID, VoiceID: VoiceID}
	players[player.GuildID] = player

	go func() {
		if session, err := NewVoicy(core.Session, GuildID, VoiceID); err == nil {
			player.Session = session
			player.Play()
		} else {
			player.Kill(true, emojis.Cry, "Um erro ocorreu ao tentar se conectar ao canal de voz: `%v`", err)
		}
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
	p.Session.Stop()
}

func (p *Player) Pause() {
	if p.Current != nil && !p.Current.IsLive {
		p.Session.Pause()
		p.State = PausedState
	}
}

func (p *Player) Resume() {
	if p.Current != nil {
		p.Session.Resume()
		p.State = PlayingState
	}
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

	if len(p.Queue) == 0 || p.Session == nil {
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

			core.SendMessage(p.TextID, emojis.Cry, "Um erro ocorreu ao carregar a música **%s**: `%v`", current.Title, err)
			return
		} else {
			current.Song = song
		}
	}

	p.Current, p.State = current, PlayingState
	go func() {
		if p.State == DestroyedState {
			return
		}

		if err := p.Session.PlayURL(current.StreamingURL, current.IsOpus); err != nil {
			core.SendMessage(p.TextID, emojis.Cry, "Um erro ocorreu enquanto tocava a música **%s**: `%v`", current.Title, err)
		}

		p.Current, p.State = nil, StoppedState
		p.Play()
	}()

	embed := core.NewEmbed().
		Description("%s Tocando agora [%s](%s)", emojis.AnimatedHype, current.Title, current.URL).
		Image(current.Thumbnail).
		Color(0x00C1FF).
		Field("Autor", current.Author, true).
		Field("Duração", utils.Is(current.IsLive, "--:--", utils.FormatTime(current.Duration)), true).
		Field("Provedor", current.Provider(), true).
		Footer(utils.Fmt("Adicionado por %s#%s", current.Requester.Username, current.Requester.Discriminator), current.Requester.AvatarURL()).
		Timestamp(current.Time)

	core.Session.SendMessage(p.TextID, "", embed.Build())
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
			p.State = DestroyedState

			if p.Timer != nil {
				p.Timer.Stop()
				p.Timer = nil
			}

			p.Queue = []*RequestedSong{}
			if p.Session != nil {
				p.Session.Destroy()
			}

			delete(players, p.GuildID)
			if len(args) >= 2 {
				core.SendMessage(p.TextID, args[0].(string), args[1].(string), args[2:]...)
			}
		}

		p.Unlock()
	}

	if force {
		p.Unlock()
		removePlayer()
		return
	}

	if p.Timer == nil && p.State == StoppedState {
		p.Timer = time.AfterFunc(5*time.Minute, removePlayer)
	}

	p.Unlock()
}
