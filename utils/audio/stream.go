package audio

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/ItsClairton/Anny/logger"
	"github.com/bwmarrin/discordgo"
)

// Isso é baseado no https://github.com/jonas747/dca porém com algumas correções e mais básico

type OpusReader interface {
	OpusFrame() (frame []byte, err error)
	FrameDuration() time.Duration
}

type StreamingSession struct {
	sync.Mutex
	source     *EncodeSession
	connection *discordgo.VoiceConnection
	running    bool
	paused     bool
	finished   bool
	framesSent int

	callback chan error
	err      error
}

func NewStream(source *EncodeSession, vc *discordgo.VoiceConnection, callback chan error) *StreamingSession {

	session := &StreamingSession{
		source:     source,
		connection: vc,
		callback:   callback,
	}

	go session.stream()

	return session
}

func (s *StreamingSession) stream() {

	s.Lock()

	if s.running {
		s.Unlock()
		logger.Warn("Voice already running")
		return
	}

	s.running = true
	s.Unlock()

	defer func() {
		s.Lock()
		s.running = false
		s.Unlock()
	}()

	for {
		s.Lock()

		if s.paused {
			s.Unlock()
			return
		}
		s.Unlock()
		err := s.readNext()

		if err != nil {
			s.Lock()
			s.finished = true

			if err != io.EOF {
				s.err = err
			}

			if s.callback != nil {
				s.callback <- err
			}

			s.Unlock()
			break
		}

	}

}

func (s *StreamingSession) readNext() error {

	opus, err := s.source.OpusFrame()
	if err != nil {
		return err
	}

	timeOut := time.After(time.Second)

	select {
	case <-timeOut:
		return errors.New("voice connection is closed")
	case s.connection.OpusSend <- opus:
	}
	s.Lock()
	s.framesSent++
	s.Unlock()

	return nil
}

func (s *StreamingSession) PlaybackPosition() int {
	s.Lock()
	time := s.framesSent * int(s.source.FrameDuration())
	s.Unlock()
	return time
}

func (s *StreamingSession) Pause() {
	s.Lock()

	if s.finished {
		s.Unlock()
		return
	}

	isPaused := !(s.paused)

	s.paused = isPaused
	if !isPaused {
		go s.stream()
	}
	s.Unlock()
}

func (s *StreamingSession) Finished() (bool, error) {
	s.Lock()

	err := s.err
	state := s.finished

	s.Unlock()
	return state, err
}

func (s *StreamingSession) Paused() bool {
	s.Lock()

	state := s.paused

	s.Unlock()
	return state
}

func (s *StreamingSession) Source() *EncodeSession {
	s.Lock()

	source := s.source

	s.Unlock()
	return source
}
