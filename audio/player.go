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
	Name, Author     string
	Requester        *discordgo.User
	IsOpus           bool
}

type CurrentTrack struct {
	*Track
	Session *StreamingSession
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

func RemovePlayer(player *Player) {
	player.connection.Disconnect()
	players[player.guild] = nil
}

func GetPlayer(id string) *Player {
	return players[id]
}

func (p *Player) Play() {
	p.Lock()

	if p.state != StoppedState {
		p.Unlock()
		return
	}

	if len(p.queue) < 1 {
		p.Unlock()
		RemovePlayer(p)
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
	p.Play()
}
