package music

import (
	"io"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/base/embed"
	"github.com/ItsClairton/Anny/base/response"
	"github.com/ItsClairton/Anny/utils/Emotes"
	"github.com/ItsClairton/Anny/utils/audio"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/bwmarrin/discordgo"
)

type Player struct {
	State      string
	Guild      *discordgo.Guild
	Connection *discordgo.VoiceConnection
	Ctx        *base.CommandContext
	Current    CurrentTrack
	Tracks     []Track
}

type Track struct {
	Name      string
	Author    string
	URL       string
	ThumbURL  string
	StreamURL string
	Duration  int64
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

func (p *Player) LoadTrack(track Track) {
	p.Tracks = append(p.Tracks, track)
	go p.Play()
}

func (p *Player) Play() {

	if p.State == PlayingState {
		return
	}

	if len(p.Tracks) < 1 {
		return
	}

	p.Current = CurrentTrack{p.Tracks[0], nil}
	p.Tracks = p.Tracks[1:]

	encodingSession := audio.EncodeData(p.Current.StreamURL)

	defer encodingSession.Cleanup()
	p.State = PlayingState

	eb := embed.NewEmbed(p.Ctx.Locale, "music.playingEmbed").
		WithEmoteDescription(Emotes.YEAH, sutils.Fmt("[%s](%s)", p.Current.Name, p.Current.URL)).
		WithField(p.Current.Author, true).
		WithField(sutils.ToHHMMSS(float64(p.Current.Duration)/1000), true).
		SetImage(p.Current.ThumbURL).
		SetColor(0x006798).
		WithFooter(p.Current.Requester.AvatarURL(""), p.Current.Requester.Username)

	p.Ctx.SendWithResponse(response.New(p.Ctx.Locale).WithEmbed(eb))

	done := make(chan error)
	p.Current.Session = audio.NewStream(encodingSession, p.Connection, done)
	err := <-done

	if err != nil {

		if err == io.EOF {
			p.State = StoppedState
			p.Play()
			return
		}
		logger.Warn(err.Error())
	}

}
