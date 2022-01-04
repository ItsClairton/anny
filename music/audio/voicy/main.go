package voicy

import (
	"bytes"
	"context"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
	"github.com/diamondburned/oggreader"
	"github.com/pkg/errors"
)

const (
	stoppedState = iota
	changingState
	pausedState
	playingState
)

type Session struct {
	Session *voice.Session

	source string
	isOpus bool

	Position time.Duration

	state   int
	channel chan int

	context context.Context
	cancel  context.CancelFunc
}

func New(state *state.State, guildID discord.GuildID, channelID discord.ChannelID) (*Session, error) {
	session, err := voice.NewSession(state)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a voice session")
	}

	if err := session.JoinChannel(guildID, channelID, false, true); err != nil {
		return nil, errors.Wrap(err, "unable to connect to voice channel")
	}

	return &Session{Session: session}, nil
}

func (s *Session) PlayURL(source string, isOpus bool) error {
	if s.state != stoppedState && s.state != changingState {
		s.Stop()
	}

	s.context, s.cancel = context.WithCancel(context.Background())
	s.source, s.isOpus = source, isOpus

	ffmpeg := exec.CommandContext(s.context, "ffmpeg",
		"-loglevel", "error", "-reconnect", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "5", "-ss", utils.FormatTime(s.Position),
		"-i", source, "-vn", "-codec", utils.Is(s.isOpus, "copy", "libopus"), "-vbr", "off", "-frame_duration", "20", "-f", "opus", "-")

	stdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		s.stop()
		return errors.Wrapf(err, "failed to get ffmpeg stdout")
	}

	var stderr bytes.Buffer
	ffmpeg.Stderr = &stderr

	if err := ffmpeg.Start(); err != nil {
		s.stop()
		return errors.Wrapf(err, "failed to start ffmpeg process")
	}

	if err := s.SendFlag(voicegateway.Microphone); err != nil {
		s.stop()
		return errors.Wrapf(err, "failed to send speaking packet to discord")
	}

	s.setState(playingState)

	if err := oggreader.DecodeBuffered(s, stdout); err != nil && s.state != changingState {
		s.stop()
		return err
	}

	if err, std := ffmpeg.Wait(), stderr.String(); err != nil && s.state != changingState && std != "" {
		s.stop()
		return errors.Wrapf(errors.New(strings.ReplaceAll(std, s.source, "source")), "ffmpeg returned error")
	}

	if s.state == changingState {
		return s.PlayURL(s.source, s.isOpus)
	}

	s.stop()
	return nil
}

func (s *Session) Destroy() {
	s.Stop()
	s.Session.Leave()
}

func (s *Session) Seek(position time.Duration) {
	if s.state == stoppedState {
		return
	}

	s.Position = position
	s.setState(changingState)
	s.Stop()
}

func (s *Session) Resume() {
	if s.state == pausedState {
		s.setState(playingState)
		s.SendFlag(voicegateway.Microphone)
	}
}

func (s *Session) Pause() {
	if s.state != stoppedState && s.state != changingState {
		s.setState(pausedState)
		s.SendFlag(voicegateway.NotSpeaking)
	}
}

func (s *Session) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Session) SendFlag(flag voicegateway.SpeakingFlag) error {
	if s.Session.VoiceUDPConn() == nil {
		return net.ErrClosed
	}

	return s.Session.Speaking(flag)
}

func (s *Session) Write(data []byte) (int, error) {
	if s.state == stoppedState || s.state == changingState {
		return 0, context.Canceled
	}

	if s.state == pausedState {
		s.channel = make(chan int)

		for {
			if newState := <-s.channel; newState != pausedState {
				close(s.channel)
				s.channel = nil
				break
			}
		}
	}

	s.Position = s.Position + (20 * time.Millisecond)
	return s.Session.WriteCtx(s.context, data)
}

func (s *Session) setState(state int) {
	s.state = state

	if s.channel != nil {
		s.channel <- state
	}
}

func (s *Session) stop() {
	s.cancel()
	s.setState(stoppedState)
	s.Position = 0
	s.SendFlag(voicegateway.NotSpeaking)
}
