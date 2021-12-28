package music

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/music/audio/voicy"
	music "github.com/ItsClairton/Anny/music/providers"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/pkg/errors"
)

const (
	StoppedState = iota
	DestroyedState
	PausedState
	PlayingState
)

var (
	players   = map[discord.GuildID]*Player{}
	TimerTime = 3 * time.Minute
)

type Player struct {
	*sync.Mutex
	*core.BasicContext

	Voicy *voicy.Session

	Timer *time.Timer
	State int

	Current *RequestedSong
	Queue   []*RequestedSong

	GuildID         discord.GuildID
	VoiceID, TextID discord.ChannelID
}

type RequestedSong struct {
	*music.Song
	Requester   *discord.User `json:"requester"`
	RequestedAt time.Time     `json:"requestedAt"`
}

func GetOrCreatePlayer(guildID discord.GuildID, textID, voiceID discord.ChannelID) *Player {
	player := players[guildID]

	if player == nil {
		player = NewPlayer(guildID, textID, voiceID)
	}

	if player.Timer != nil {
		player.Timer.Stop()
	}

	return player
}

func GetPlayer(guildID discord.GuildID) *Player {
	return players[guildID]
}

func NewPlayer(guildID discord.GuildID, textID, voiceID discord.ChannelID) *Player {
	if players[guildID] != nil {
		logger.Warn(errors.New("something tried to create a new player for a guild that already has an existing player"))
		return players[guildID]
	}

	player := &Player{
		Mutex:        &sync.Mutex{},
		BasicContext: core.NewBasicContext(textID, guildID),
		State:        StoppedState,
		GuildID:      guildID,
		TextID:       textID,
		VoiceID:      voiceID,
	}

	players[guildID] = player
	go func() {
		if session, err := voicy.New(core.State, guildID, voiceID); err == nil {
			player.Voicy = session
			player.Play()
		} else {
			player.Send(emojis.Cry, "Um erro ocorreu ao tentar se conectar ao canal de voz: ```go\n%+v```", err)
			player.Stop(false)
		}
	}()

	return player
}

func (p *Player) Play() {
	if p.State != StoppedState || p.Voicy == nil {
		return
	}

	if len(p.Queue) == 0 {
		p.Stop(true)
		return
	}

	if p.Timer != nil {
		p.Timer.Stop()
	}

	defer p.Play()

	song := p.Queue[0]
	if err := song.Load(); err != nil {
		p.Queue = p.Queue[1:]
		p.Send(emojis.Cry, "Um erro ocorreu ao carregar a música **%s**: `%v`", song.Title, err)
		return
	}

	p.Queue, p.Current, p.State = p.Queue[1:], song, PlayingState
	go p.Send(utils.NewEmbed().
		Description("%s Tocando agora [%s](%s)", emojis.AnimatedHype, song.Title, song.URL).
		Image(song.Thumbnail).
		Color(0x00C1FF).
		Field("Autor", song.Author, true).
		Field("Duração", utils.Is(song.IsLive, "--:--", utils.FormatTime(song.Duration)), true).
		Field("Provedor", song.Provider(), true).
		Timestamp(song.RequestedAt).
		Footer(utils.Fmt("Adicionado por %s#%s", song.Requester.Username, song.Requester.Discriminator), song.Requester.AvatarURL()))

	if err := p.Voicy.PlayURL(song.StreamingURL, song.IsOpus); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		p.Send(emojis.Cry, "Um erro ocorreu ao tocar a música **%s**: `%v`", song.Title, err)
	}

	p.Current, p.State = nil, StoppedState
}

func (p *Player) Stop(schedule bool) {
	p.Lock()
	defer p.Unlock()

	if !schedule {
		removePlayer(p, false)
		return
	}

	if p.State != StoppedState || len(p.Queue) != 0 {
		logger.Warn(errors.New("something tried to start the timer, but the player is still playing something"))
		return
	}

	if p.Timer != nil {
		p.Timer.Stop()
	}

	p.Timer = time.AfterFunc(TimerTime, func() {
		p.Lock()
		defer p.Unlock()

		removePlayer(p, true)
	})
}

func (p *Player) Pause() {
	if p.State == PlayingState {
		p.Voicy.Pause()
		p.State = PausedState
	}
}

func (p *Player) Resume() {
	if p.State == PausedState {
		p.Voicy.Resume()
		p.State = PlayingState
	}
}

func (p *Player) Skip() {
	if p.Current == nil {
		return
	}

	p.Current, p.State = nil, StoppedState
	p.Voicy.Stop()
}

func (p *Player) AddSong(requester *discord.User, shuffle bool, songs ...*music.Song) {
	for _, song := range songs {
		p.Queue = append(p.Queue, &RequestedSong{song, requester, time.Now()})
	}
	if shuffle {
		p.Shuffle()
	}

	go p.Play()
}

func (p *Player) Shuffle() {
	rand.Shuffle(len(p.Queue), func(oldPos, newPos int) {
		p.Queue[oldPos], p.Queue[newPos] = p.Queue[newPos], p.Queue[oldPos]
	})
}

func removePlayer(player *Player, scheduled bool) {
	if player == nil || player.State == DestroyedState {
		return
	}

	if scheduled && (player.State != StoppedState || len(player.Queue) != 0) {
		return
	}

	player.State = DestroyedState
	player.Queue = []*RequestedSong{}

	if player.Timer != nil {
		player.Timer.Stop()
	}

	if player.Voicy != nil {
		player.Voicy.Destroy()
	}

	delete(players, player.GuildID)
}
