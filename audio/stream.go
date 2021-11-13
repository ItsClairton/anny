package audio

import (
	"io"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
)

// https://github.com/jonas747/dca

type OpusReader interface {
	OpusFrame() (frame []byte, err error)
	FrameDuration() time.Duration
}

type StreamingSession struct {
	sync.Mutex
	source     *EncodingSession
	connection *voice.Session

	running  bool
	paused   bool
	finished bool

	framesSent int
	callback   chan error
}

func NewStream(source *EncodingSession, vc *voice.Session, callback chan error) *StreamingSession {
	session := &StreamingSession{
		source:     source,
		connection: vc,
		callback:   callback,
	}
	go session.stream()

	return session
}

func StreamURL(URL string, connection *voice.Session, callback chan error) *StreamingSession {
	return NewStream(NewEncodingURL(URL), connection, callback)
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

	if err := s.connection.Speaking(voicegateway.Microphone); err != nil {
		s.callback <- err
		close(s.callback)
		return
	}

	for {
		s.Lock()

		if s.finished {
			s.Unlock()
			break
		}

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
			s.callback <- err

			close(s.callback)
			s.source.Stop()
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

	_, err = s.connection.Write(opus)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	s.framesSent++
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

func (s *StreamingSession) Stop() {
	s.Lock()
	defer s.Unlock()

	s.finished = true
	s.source.Stop()

	if s.source.err != nil {
		s.callback <- s.source.err
		return
	}

	s.callback <- io.EOF
}
