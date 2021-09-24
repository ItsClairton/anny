package audio

import (
	"io"
	"sync"

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
	sync.Mutex
	state      int
	connection *discordgo.VoiceConnection
	queue      []*Track
	current    *CurrentTrack

	guild string
}

type Track struct {
	ID, StreamingUrl string
	Title, Author    string
	Requester        *discordgo.User
	IsOpus           bool
}

type CurrentTrack struct {
	*Track
	Session *StreamingSession
}

func GetOrCreatePlayer(s *discordgo.Session, guildId, voiceId string) (*Player, error) {
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

	player := NewPlayer(guildId, conn)
	AddPlayer(player)
	return player, nil
}

func NewPlayer(guild string, conn *discordgo.VoiceConnection) *Player {
	return &Player{
		state:      StoppedState,
		connection: conn,
		guild:      guild,
	}
}

func AddPlayer(player *Player) *Player {
	players[player.guild] = player
	return player
}

func RemovePlayer(player *Player, force bool) {
	player.Lock()
	if !force && (player.state != StoppedState || len(player.queue) > 0) {
		player.Unlock()
		return
	}

	player.connection.Disconnect()
	players[player.guild] = nil
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

	session := NewProcessingSession(p.current.StreamingUrl, p.current.IsOpus)
	defer session.StopClean()

	done := make(chan error)
	p.current.Session = NewStream(session, p.connection, done)
	p.state = PlayingState
	p.Unlock()
	err := <-done
	if err != nil {
		if err != io.EOF {
			logger.Warn(err.Error())
		}
		p.Lock()
		p.state = StoppedState
		p.Unlock()
		p.Play()
	}
}

func (p *Player) AddQueue(track *Track) {
	p.Lock()
	p.queue = append(p.queue, track)
	p.Unlock()
	go p.Play()
}

func (p *Player) GetQueue() []*Track {
	p.Lock()
	queue := p.queue
	p.Unlock()

	return queue
}

func (p *Player) GetState() int {
	p.Lock()
	state := p.state
	p.Unlock()

	return state
}
