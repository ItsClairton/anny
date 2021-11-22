package audio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
	"github.com/diamondburned/oggreader"
)

var stoppedState, changeState, pausedState, playingState = 0, 1, 2, 3

type VoicySession struct {
	Session *voice.Session

	source string
	isOpus bool
	volume float32

	Position time.Duration

	state   int
	channel chan int

	context context.Context
	cancel  context.CancelFunc
}

func NewVoicy(state *state.State, guildID discord.GuildID, channelID discord.ChannelID) (*VoicySession, error) {
	voice, err := voice.NewSession(state)
	if err != nil {
		return nil, err
	}

	if err := voice.JoinChannel(guildID, channelID, false, true); err != nil {
		return nil, err
	}

	return &VoicySession{Session: voice, volume: 1}, nil
}

func (vs *VoicySession) SetVolume(volume int64) {
	vol := float32(volume) / 100
	if vol == vs.volume {
		return
	}

	if vol < 1 {
		vol = -vol
	}

	vs.volume = vol
	vs.setState(changeState)
	vs.Stop()
}

func (vs *VoicySession) SetPosition(duration time.Duration) {
	vs.Position = duration

	vs.setState(changeState)
	vs.Stop()
}

func (vs *VoicySession) PlayURL(source string, isOpus bool) error {
	if vs.state != stoppedState && vs.state != changeState {
		vs.Stop()
	}

	vs.context, vs.cancel = context.WithCancel(context.Background())
	vs.source, vs.isOpus = source, isOpus
	defer vs.Stop()

	options := []string{"-reconnect", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "5", "-ss", utils.FormatTime(vs.Position), "-loglevel", "error", "-i", source, "-vn", "-vbr", "on", "-frame_duration", "20"}

	if isOpus && vs.volume == 1 {
		options = append(options, "-preset", "ultrafast", "-codec", "copy")
	}

	if vs.volume != 1 {
		options = append(options, "-af", fmt.Sprintf("volume=%v", vs.volume))
	}

	ffmpeg := exec.CommandContext(vs.context, "ffmpeg", append(options, "-f", "opus", "-")...)
	stdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get ffmpeg stdout: %w", err)
	}

	var stderr bytes.Buffer
	ffmpeg.Stderr = &stderr

	if err := vs.SendSpeaking(); err != nil {
		return fmt.Errorf("failed to send speaking packet to discord: %w", err)
	}

	if err := ffmpeg.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg process: %w", err)
	}

	vs.setState(playingState)

	if err := oggreader.DecodeBuffered(vs, stdout); err != nil && err.Error() != "failed to write a packet: context canceled" {
		return fmt.Errorf("failed to send to voice connection: %w", err)
	}

	if err := ffmpeg.Wait(); err != nil && err.Error() != "signal: killed" {
		if str := string(stderr.String()); str != "" {
			return fmt.Errorf("ffmpeg returned error: %w", errors.New(strings.ReplaceAll(str, vs.source, "source")))
		}
		return fmt.Errorf("failed to wait ffmpeg: %w", err)
	}

	if vs.state == changeState {
		return vs.PlayURL(vs.source, vs.isOpus)
	}

	return nil
}

func (vs *VoicySession) Destroy() {
	vs.Stop()
	vs.Session.Leave()
}

func (vs *VoicySession) State() int {
	return vs.state
}

func (vs *VoicySession) Resume() {
	if vs.state == pausedState {
		vs.setState(playingState)
		vs.SendSpeaking()
	}
}

func (vs *VoicySession) Pause() {
	if vs.state != stoppedState {
		vs.setState(pausedState)
	}
}

func (vs *VoicySession) Stop() {
	if vs.state != changeState {
		vs.Position = 0
		vs.setState(stoppedState)
	}

	if vs.cancel != nil {
		vs.cancel()
	}
}

func (vs *VoicySession) SendSpeaking() error {
	return vs.Session.Speaking(voicegateway.Microphone)
}

func (vs *VoicySession) Write(data []byte) (n int, err error) {
	if vs.state == pausedState {
		vs.waitState(playingState, stoppedState)
	}

	if vs.state == stoppedState {
		return 0, nil
	}

	vs.Position = vs.Position + (20 * time.Millisecond)
	return vs.Session.WriteCtx(vs.context, data)
}

func (vs *VoicySession) waitState(states ...int) {
	vs.channel = make(chan int)

	for {
		if newState := <-vs.channel; utils.IntegerArrayContains(states, newState) {
			close(vs.channel)
			vs.channel = nil
			break
		}
	}
}

func (vs *VoicySession) setState(state int) {
	vs.state = state

	if vs.channel != nil {
		vs.channel <- state
	}
}
