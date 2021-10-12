package audio

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// https://github.com/jonas747/dca

type OpusReader interface {
	OpusFrame() (frame []byte, err error)
	FrameDuration() time.Duration
}

type StreamingSession struct {
	sync.Mutex
	source     *EncodingSession
	connection *discordgo.VoiceConnection

	running  bool
	paused   bool
	finished bool

	framesSent int
	callback   chan error
}

func NewStream(source *EncodingSession, vc *discordgo.VoiceConnection, callback chan error) *StreamingSession {
	session := &StreamingSession{
		source:     source,
		connection: vc,
		callback:   callback,
	}
	go session.stream()

	return session
}

func StreamFromPath(path string, connection *discordgo.VoiceConnection, callback chan error) *StreamingSession {
	return NewStream(NewEncodingFromPath(path), connection, callback)
}

func StreamFromReader(reader io.Reader, connection *discordgo.VoiceConnection, callback chan error) *StreamingSession {
	return NewStream(NewEncodingFromReader(reader), connection, callback)
}

func (s *StreamingSession) stream() {
	s.Lock()

	if s.running {
		s.Unlock()
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
			if s.source.err != nil {
				err = s.source.err
			}

			if s.callback != nil {
				s.callback <- err
			}
			s.source.StopClean()
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

	timeOut := time.After(40 * time.Second)
	select {
	case <-timeOut:
		return errors.New("voice connection timeout")
	case s.connection.OpusSend <- opus:
	}
	s.Lock()
	s.framesSent++
	s.Unlock()

	return nil
}

func (s *StreamingSession) PlaybackPosition() time.Duration {
	s.Lock()
	defer s.Unlock()

	return time.Duration(s.framesSent) * s.source.FrameDuration()
}

func (s *StreamingSession) Pause(paused bool) {
	s.Lock()
	defer s.Unlock()

	if s.finished {
		return
	}

	s.paused = paused
	if !paused {
		go s.stream()
	}
}

func (s *StreamingSession) Finished() bool {
	s.Lock()
	defer s.Unlock()

	return s.finished
}

func (s *StreamingSession) Paused() bool {
	s.Lock()
	defer s.Unlock()

	return s.paused
}

func (s *StreamingSession) Source() *EncodingSession {
	s.Lock()
	defer s.Unlock()

	return s.source
}
