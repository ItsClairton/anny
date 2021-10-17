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
	state      int
	connection *discordgo.VoiceConnection
	queue      []*RequestedSong
	current    *CurrentSong

	guildID, textID, voiceID string
}

type RequestedSong struct {
	*Song
	Requester *discordgo.User
	Time      time.Time
}

type CurrentSong struct {
	*RequestedSong
	Session *StreamingSession
}

func NewPlayer(guildID, textID, voiceID string, conn *discordgo.VoiceConnection) *Player {
	return &Player{
		Mutex:      &sync.Mutex{},
		state:      StoppedState,
		connection: conn,
		guildID:    guildID,
		textID:     textID,
		voiceID:    voiceID,
	}
}

func GetPlayer(id string) *Player {
	return players[id]
}

func AddPlayer(player *Player) *Player {
	players[player.guildID] = player
	return player
}

func RemovePlayer(player *Player, force bool) {
	removeFunc := func() {
		player.Lock()
		defer player.Unlock()
		if player.connection != nil {
			player.connection.Disconnect()
		}
		players[player.guildID] = nil
	}

	if force {
		removeFunc()
	} else {
		if player.timer == nil {
			player.timer = time.AfterFunc(5*time.Minute, removeFunc)
		}
	}
}

func (p *Player) UpdateVoice(voiceID string, connection *discordgo.VoiceConnection) {
	p.Lock()
	defer p.Unlock()

	p.connection, p.voiceID = connection, voiceID
	go p.Play()
}

func (p *Player) AddSong(requester *discordgo.User, tracks ...*Song) {
	p.Lock()
	for _, track := range tracks {
		p.queue = append(p.queue, &RequestedSong{track, requester, time.Now()})
	}
	p.Unlock()
	go p.Play()
}

func (p *Player) GetQueue() []*RequestedSong {
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

func (p *Player) GetCurrent() *CurrentSong {
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

	if p.state != StoppedState || p.connection == nil {
		p.Unlock()
		return
	}

	if len(p.queue) < 1 {
		p.Unlock()
		RemovePlayer(p, false)
		return
	}

	if p.timer != nil {
		p.timer.Stop()
		p.timer = nil
	}

	p.current, p.queue = &CurrentSong{p.queue[0], nil}, p.queue[1:]
	current := p.current

	if current.StreamingURL == "" {
		song, err := current.Provider.GetInfo(current.Song)
		if err != nil {
			p.sendError(err)
			p.Unlock()
			go p.Play()
			return
		}
		p.current.RequestedSong.Song = song
		current = p.current
	}

	done := make(chan error)
	current.Session, p.state = StreamFromPath(current.StreamingURL, p.connection, done), PlayingState
	p.Unlock()

	go func() {
		embed := discord.NewEmbed().
			SetDescription(utils.Fmt("%s Tocando agora [%s](%s)", emojis.ZeroYeah, current.Title, current.URL)).
			SetThumbnail(current.Thumbnail).
			SetColor(0xA652BB).
			AddField("Autor", current.Author, true).
			AddField("Duração", utils.ToDisplayTime(current.Duration.Seconds()), true).
			AddField("Provedor", current.Provider.PrettyName(), true).
			SetTimestamp(current.Time.Format(time.RFC3339))
		if current.Playlist != nil {
			embed.SetFooter(utils.Fmt("Pedido por %s • Playlist %s", current.Requester.Username, current.Playlist.Title), current.Requester.AvatarURL(""))
		} else {
			embed.SetFooter(utils.Fmt("Pedido por %s", current.Requester.Username), current.Requester.AvatarURL(""))
		}

		discord.NewResponse().WithEmbed(embed).Send(p.textID)
	}()

	err := <-done
	if err != nil {
		p.Lock()
		defer p.Unlock()
		if err != io.EOF {
			p.sendError(err)
		}

		p.state = StoppedState
		go p.Play()
	}
}

func (p *Player) sendError(err error) {
	discord.NewResponse().WithEmbed(
		discord.NewEmbed().
			SetColor(0xF93A2F).
			SetDescription(utils.Fmt("%s Um erro ocorreu ao tocar [%s](%s): `%v`", emojis.MikuCry, p.current.Title, p.current.URL, err)),
	).Send(p.textID)
}
