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
	sync.Mutex
	state      int
	connection *discordgo.VoiceConnection
	queue      []*Track
	current    *CurrentTrack

	guildId, textId, voiceId string
}

type Track struct {
	URL, ID       string
	Title, Author string
	Requester     *discordgo.User
	Duration      time.Duration
	IsOpus        bool
	Playlist      *Playlist

	StreamingUrl, ThumbnailUrl string
}

type Playlist struct {
	ID, Title, Author string
}

type CurrentTrack struct {
	*Track
	Session *StreamingSession
}

func GetOrCreatePlayer(s *discordgo.Session, guildId, textId, voiceId string) (*Player, error) {
	if GetPlayer(guildId) != nil {
		return GetPlayer(guildId), nil
	}

	conn, err := s.ChannelVoiceJoin(guildId, voiceId, false, true)
	if err != nil {
		conn, err := s.ChannelVoiceJoin(guildId, voiceId, false, true)

		if err != nil {
			if conn != nil {
				conn.Disconnect()
			}
			return nil, err
		}
	}

	player := NewPlayer(guildId, textId, voiceId, conn)
	AddPlayer(player)
	return player, nil
}

func NewPlayer(guildId, textId, voiceId string, conn *discordgo.VoiceConnection) *Player {
	return &Player{
		state:      StoppedState,
		connection: conn,
		guildId:    guildId,
		textId:     textId,
		voiceId:    voiceId,
	}
}

func AddPlayer(player *Player) *Player {
	players[player.guildId] = player
	return player
}

func RemovePlayer(player *Player, force bool) {
	player.Lock()
	if !force && (player.state != StoppedState || len(player.queue) > 0) {
		player.Unlock()
		return
	}

	player.connection.Disconnect()
	players[player.guildId] = nil
	player.Unlock()
}

func GetPlayer(id string) *Player {
	return players[id]
}

func (p *Player) Skip() {
	p.Lock()
	p.current.Session.source.StopClean()
	p.Unlock()
}

func (p *Player) Pause() {
	p.Lock()
	p.current.Session.Pause(true)
	p.state = PausedState
	p.Unlock()
}

func (p *Player) Unpause() {
	p.Lock()
	p.current.Session.Pause(false)
	p.state = PlayingState
	p.Unlock()
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

	if p.current.StreamingUrl == "" {
		track, err := GetTrack(p.current.ID, p.current.Requester)
		if err != nil {
			discord.NewResponse().
				WithContentEmoji(emojis.MikuCry, "Um erro ocorreu ao tocar a música **%s**: `%s`", p.current.Title, err.Error()).
				SendTo(p.textId)
			p.Unlock()
			RemovePlayer(p, false)
			return
		}
		p.current.Track = track
	}

	done := make(chan error)
	p.current.Session = NewStream(NewProcessingSession(p.current.StreamingUrl, p.current.IsOpus), p.connection, done)
	p.state = PlayingState
	p.Unlock()

	discord.NewResponse().
		WithEmbed(discord.NewEmbed().
			SetDescription(utils.Fmt("%s Tocando agora [%s](%s)", emojis.ZeroYeah, p.current.Title, p.current.URL)).
			SetThumbnail(p.current.ThumbnailUrl).
			SetColor(0xA652BB).
			AddField("Autor", p.current.Author, true).
			AddField("Duração", utils.ToDisplayTime(p.current.Duration.Seconds()), true).
			SetFooter(utils.Fmt("Pedido por %s", p.current.Requester.Username), p.current.Requester.AvatarURL("")).
			Build()).SendTo(p.textId)

	go func(p *Player) {
		p.Lock()
		if len(p.queue) > 0 && p.queue[0].StreamingUrl == "" {
			track, err := GetTrack(p.queue[0].ID, p.queue[0].Requester)
			if err == nil {
				p.queue[0] = track
			}
		}
		p.Unlock()
	}(p)

	err := <-done
	if err != nil {
		p.Lock()
		if err != io.EOF {
			discord.NewResponse().
				WithContentEmoji(emojis.MikuCry, "Um erro ocorreu ao tocar a música **%s**: `%s`", p.current.Title, err.Error()).
				SendTo(p.textId)
		}
		p.state = StoppedState
		p.Unlock()
		p.Play()
	}
}

func (p *Player) AddTrack(track *Track) {
	p.Lock()
	p.queue = append(p.queue, track)
	p.Unlock()
	go p.Play()
}

func (p *Player) AddPlaylist(tracks []*Track) {
	p.Lock()
	p.queue = append(p.queue, tracks...)
	p.Unlock()
	go p.Play()
}

func (p *Player) GetQueue() []*Track {
	p.Lock()
	queue := p.queue
	p.Unlock()

	return queue
}

func (p *Player) Shuffle() {
	p.Lock()
	rand.Shuffle(len(p.queue), func(old, new int) {
		p.queue[old], p.queue[new] = p.queue[new], p.queue[old]
		p.queue[new].StreamingUrl = ""
		p.queue[old].StreamingUrl = ""
	})
	p.Unlock()
}

func (p *Player) GetCurrent() *CurrentTrack {
	p.Lock()
	current := p.current
	p.Unlock()

	return current
}

func (p *Player) GetState() int {
	p.Lock()
	state := p.state
	p.Unlock()

	return state
}
