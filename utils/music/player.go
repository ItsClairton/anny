package music

import (
	"io"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/embed"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/audio"
	"github.com/ItsClairton/Anny/utils/constants"
	"github.com/ItsClairton/Anny/utils/music/provider"
	"github.com/bwmarrin/discordgo"
)

type Player struct {
	State       string
	GuildID     string
	Connection  *discordgo.VoiceConnection
	Ctx         *base.CommandContext
	Tracks      []Track
	Current     CurrentTrack
	lastMessage string
}

type Track struct {
	*provider.PartialInfo
	Stream    *provider.StreamInfo
	Requester *discordgo.User
}

type CurrentTrack struct {
	Track
	Session *audio.StreamingSession
}

var (
	PlayingState = "PLAYING"
	StoppedState = "STOPPED"
	PausedState  = "PAUSED"
)

func (p *Player) AddQueue(track Track) {
	p.Tracks = append(p.Tracks, track)
	go p.Play()
	go p.loadNextTrack()
}

func (p *Player) loadNextTrack() {
	if len(p.Tracks) < 1 {
		return
	}

	track := &p.Tracks[0]
	if track.Stream != nil {
		return
	}

	stream, err := track.Provider.GetStream(track.PartialInfo)

	if err != nil {
		return
	}
	track.Stream = stream
}

func (p *Player) Play() {

	if p.State != StoppedState {
		return
	}

	if len(p.Tracks) < 1 {
		RemovePlayer(p)
		return
	}

	p.Current = CurrentTrack{p.Tracks[0], nil}
	p.Tracks = p.Tracks[1:]

	track := &p.Current

	if track.Stream == nil {
		stream, err := track.Provider.GetStream(track.PartialInfo)

		if err != nil {
			p.Ctx.Reply(constants.MIKU_CRY, "music.error", track.Title, track.Requester.Mention(), err.Error())
			return
		}
		track.Stream = stream
	}
	session := audio.EncodeData(track.Stream.StreamURL, track.Stream.IsOpus)
	defer session.Cleanup()

	done := make(chan error)
	track.Session = audio.NewStream(session, p.Connection, done)

	p.State = PlayingState
	eb := embed.NewEmbed(p.Ctx.Locale, "music.playingEmbed").
		WithEmoteDescription(constants.YEAH, utils.Fmt("[%s](%s)", track.Title, track.URL)).
		WithField(track.Author, true).
		WithField(track.Duration, true).
		SetImage(track.ThumbURL).
		SetColor(0x006798).
		WithFooter(track.Requester.AvatarURL(""), track.Requester.Username)

	msg, err := p.Ctx.SendWithResponse(response.New(p.Ctx.Locale).WithEmbed(eb))
	if err == nil {
		p.lastMessage = msg.ID
	}
	go p.loadNextTrack()

	err = <-done
	if err != nil {

		if err != io.EOF {
			p.Ctx.Reply(constants.MIKU_CRY, "music.error", track.Title, track.Requester.Mention(), err.Error())
		}

		p.State = StoppedState
		p.Play()
	}

}
