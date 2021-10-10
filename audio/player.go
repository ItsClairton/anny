package audio

import (
	"io"
	"math/rand"
	"sync"

	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/providers"
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
	state      int
	connection *discordgo.VoiceConnection
	queue      []*Track
	current    *CurrentTrack

	guildId, textId, voiceId string
}

type Track struct {
	*providers.Song
	Requester *discordgo.User
}

type CurrentTrack struct {
	*Track
	Session *StreamingSession
}

func NewPlayer(guildId, textId, voiceId string, conn *discordgo.VoiceConnection) *Player {
	return &Player{
		Mutex:      &sync.Mutex{},
		state:      StoppedState,
		connection: conn,
		guildId:    guildId,
		textId:     textId,
		voiceId:    voiceId,
	}
}

func GetPlayer(id string) *Player {
	return players[id]
}

func AddPlayer(player *Player) *Player {
	players[player.guildId] = player
	return player
}

func RemovePlayer(player *Player, force bool) {
	player.Lock()
	defer player.Unlock()
	if !force && (player.state != StoppedState || len(player.queue) > 0) {
		player.Unlock()
		return
	}

	player.connection.Disconnect()
	players[player.guildId] = nil
}

func (p *Player) AddTrack(tracks ...*Track) {
	p.Lock()
	defer p.Unlock()

	p.queue = append(p.queue, tracks...)
	go p.Play()
}

func (p *Player) GetQueue() []*Track {
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

func (p *Player) GetCurrent() *CurrentTrack {
	p.Lock()
	defer p.Unlock()

	return p.current
}

func (p *Player) GetState() int {
	p.Lock()
	defer p.Unlock()

	return p.state
}

func (p *Player) Skip() {
	p.Lock()
	defer p.Unlock()
	p.current.Session.source.StopClean()
}

func (p *Player) Pause() {
	p.Lock()
	defer p.Unlock()
	p.current.Session.Pause(true)
	p.state = PausedState
}

func (p *Player) Unpause() {
	p.Lock()
	defer p.Unlock()
	p.current.Session.Pause(false)
	p.state = PlayingState
}

func (p *Player) Play() {
	p.Lock()

	if p.state != StoppedState {
		p.Unlock()
		return
	}
	if len(p.queue) < 1 {
		p.Unlock()
		RemovePlayer(p, false)
		return
	}

	p.current = &CurrentTrack{p.queue[0], nil}
	p.queue = p.queue[1:]
	current := p.current

	done := make(chan error)
	p.current.Session = NewStream(NewProcessingSession(current.DirectURL, current.IsOpus), p.connection, done)
	p.state = PlayingState
	p.Unlock()

	discord.NewResponse().
		WithEmbed(discord.NewEmbed().
			SetDescription(utils.Fmt("%s Tocando agora [%s](%s)", emojis.ZeroYeah, current.Title, current.PageURL)).
			SetThumbnail(current.ThumbnailURL).
			SetColor(0xA652BB).
			AddField("Autor", current.Uploader, true).
			AddField("Duração", current.Duration, true).
			AddField("Provedor", current.DisplayProvider(), true).
			SetFooter(utils.Fmt("Pedido por %s", current.Requester.Username), current.Requester.AvatarURL("")).
			Build()).SendTo(p.textId)

	err := <-done
	if err != nil {
		p.Lock()
		if err != io.EOF {
			discord.NewResponse().
				WithContentEmoji(emojis.MikuCry, "Um erro ocorreu ao tocar a música **%s**: `%s`", current.Title, err.Error()).
				SendTo(p.textId)
		}
		p.state = StoppedState
		p.Unlock()
		p.Play()
	}
}
