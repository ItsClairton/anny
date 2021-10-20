package audio

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/bwmarrin/discordgo"
)

var (
	StoppedState = 0
	PausedState  = 1
	PlayingState = 2

	players = map[string]*Player{}
)

type Player struct {
	*sync.Mutex

	timer      *time.Timer
	connection *discordgo.VoiceConnection

	state   int
	current *CurrentSong
	queue   []*RequestedSong

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

func GetPlayer(id string) *Player {
	return players[id]
}

func NewPlayer(guildID, textID, voiceID string) *Player {
	player := &Player{
		Mutex:   &sync.Mutex{},
		state:   StoppedState,
		GuildID: guildID, TextID: textID, VoiceID: voiceID}

	players[player.GuildID] = player
	go player.CheckConnection()
	return player
}

func RemovePlayer(player *Player) {
	player.Lock()
	defer player.Unlock()

	player.queue = []*RequestedSong{}
	if player.current != nil && player.current.StreamingSession != nil {
		player.current.source.StopClean()
	}
	if player.connection != nil {
		player.connection.Disconnect()
	}

	players[player.GuildID] = nil
}

func (p *Player) CheckConnection() {
	if p.connection == nil {
		connection, err := discord.Session.ChannelVoiceJoin(p.GuildID, p.VoiceID, false, true)
		if err != nil {
			discord.NewResponse().WithContent(emojis.MikuCry, "Um erro ocorreu na conexão com o Canal de voz.").Send(p.TextID)
			p.Kill(true)
		} else {
			p.connection = connection
			go p.Play()
		}
	}
}

func (p *Player) Current() *CurrentSong {
	p.Lock()
	defer p.Unlock()

	return p.current
}

func (p *Player) State() int {
	p.Lock()
	defer p.Unlock()

	return p.state
}

func (p *Player) Skip() {
	p.Lock()
	defer p.Unlock()
	p.current.source.StopClean()
}

func (p *Player) Pause() {
	p.Lock()
	defer p.Unlock()

	p.current.Pause(true)
	p.state = PausedState
}

func (p *Player) Unpause() {
	p.Lock()
	defer p.Unlock()

	p.current.Pause(false)
	p.state = PlayingState
}

func (p *Player) AddSong(requester *discordgo.User, tracks ...*Song) {
	p.Lock()
	defer p.Unlock()

	for _, track := range tracks {
		p.queue = append(p.queue, &RequestedSong{track, requester, time.Now()})
	}
	go p.Play()
}

func (p *Player) Queue() []*RequestedSong {
	p.Lock()
	defer p.Unlock()

	return p.queue
}

func (p *Player) Shuffle() {
	p.Lock()
	defer p.Unlock()

	rand.Shuffle(len(p.queue), func(old, new int) {
		p.queue[old], p.queue[new] = p.queue[new], p.queue[old]
	})
}

func (p *Player) Play() {
	p.Lock()
	defer p.Unlock()

	if p.state != StoppedState {
		return
	}

	if len(p.queue) == 0 || p.connection == nil {
		go p.Kill(false)
		return
	}

	if p.timer != nil {
		p.timer.Stop()
		p.timer = nil
	}

	current := p.queue[0]
	if current.StreamingURL == "" {
		song, err := current.Provider.GetInfo(current.Song)
		if err != nil {
			p.sendError(err)
			go p.Play()
			return
		}

		current.Song = song
	}

	done := make(chan error)
	p.current, p.queue, p.state = &CurrentSong{
		RequestedSong:    current,
		StreamingSession: StreamFromPath(current.StreamingURL, p.connection, done),
	}, p.queue[1:], PlayingState

	go func() {
		finished := <-done

		p.Lock()
		defer p.Unlock()
		if finished != io.EOF {
			p.sendError(finished)
		}

		p.state = StoppedState
		go p.Play()
	}()

	embed := discord.NewEmbed().
		SetDescription("%s Tocando agora: [%s](%s)", emojis.ZeroYeah, current.Title, current.URL).
		SetImage(current.Thumbnail).
		SetColor(0xA652BB).
		AddField("Autor", current.Author, true).
		AddField("Duração", utils.Is(current.IsLive, "--:--", utils.FormatTime(current.Duration)), true).
		AddField("Provedor", current.Provider.Name(), true).
		SetFooter(utils.Fmt("Pedido por %s", current.Requester.Username), current.Requester.AvatarURL("")).
		SetTimestamp(current.Time.Format(time.RFC3339))

	discord.NewResponse().WithEmbed(embed).Send(p.TextID)
}

func (p *Player) Kill(force bool) {
	p.Lock()
	defer p.Unlock()

	if force {
		go RemovePlayer(p)
		return
	}

	if p.timer == nil && p.state == StoppedState && len(p.queue) == 0 {
		p.timer = time.AfterFunc(5*time.Minute, func() { RemovePlayer(p) })
	}
}

func (p *Player) sendError(err error) {
	discord.NewResponse().
		WithContent(emojis.MikuCry, "Um erro ocorreu ao tocar a música %s: `%v`", p.current.Title, err).Send(p.TextID)
}
