package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/jonas747/ogg"
)

// https://github.com/jonas747/dca

type EncodingSession struct {
	sync.Mutex

	path string

	running   bool
	process   *os.Process
	lastFrame int

	err    error
	stderr bytes.Buffer

	data   chan []byte
	buffer bytes.Buffer
}

func NewEncodingURL(path string) *EncodingSession {
	session := &EncodingSession{path: path, data: make(chan []byte)}
	go session.start()
	return session
}

func (s *EncodingSession) start() {
	defer func() {
		s.Lock()
		close(s.data)
		s.running = false
		s.Unlock()
	}()

	s.Lock()
	s.running = true

	arguments := []string{
		"-hide_banner", "-threads", "1", "-loglevel", "error",
		"-reconnect", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "5",
		"-i", s.path, "-vn", "-c:a", "libopus", "-b:a", "96k", "-frame_duration", "20", "-vbr", "off",
		"-f", "ogg", "-"}

	cmd := exec.Command("ffmpeg", arguments...)
	cmd.Stderr = &s.stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s.err = err
		s.Unlock()
		return
	}

	err = cmd.Start()
	if err != nil {
		s.err = err
		s.Unlock()
		return
	}

	s.process = cmd.Process
	s.Unlock()
	defer func() {
		if s.err == nil {
			stderr := strings.TrimSpace(s.stderr.String())
			if stderr != "" && stderr != "<nil>" && !strings.Contains(stderr, "Error in the pull function") {
				s.Lock()
				defer s.Unlock()

				s.err = errors.New(strings.ReplaceAll(stderr, s.path, "source"))
			}
		}
	}()

	s.readStdout(stdout)
	cmd.Wait()
}

func (s *EncodingSession) Stop() {
	s.Lock()

	if s.running && s.process != nil {
		s.process.Kill()
	}
	s.Unlock()

	for range s.data {
	}
}

func (s *EncodingSession) FrameDuration() time.Duration {
	return time.Duration(20) * time.Millisecond
}

func (s *EncodingSession) OpusFrame() ([]byte, error) {
	frame := <-s.data
	if frame == nil {
		return nil, io.EOF
	}

	if len(frame) < 2 {
		return nil, errors.New("bad opus frame")
	}
	return frame[2:], nil
}

func (s *EncodingSession) ReadFrame() ([]byte, error) {
	frame := <-s.data
	if frame == nil {
		return nil, io.EOF
	}

	return frame, nil
}

func (s *EncodingSession) Read(p []byte) (int, error) {
	if s.buffer.Len() >= len(p) {
		return s.buffer.Read(p)
	}

	for s.buffer.Len() < len(p) {
		frame, err := s.ReadFrame()
		if err != nil {
			break
		}
		s.buffer.Write(frame)
	}

	return s.buffer.Read(p)
}

func (ps *EncodingSession) readStdout(std io.ReadCloser) {
	decoder := ogg.NewPacketDecoder(ogg.NewDecoder(std))
	skipPackets := 2

	for {
		packet, _, err := decoder.Decode()
		if skipPackets > 0 {
			skipPackets--
			continue
		}

		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				ps.Lock()
				ps.err = err
				ps.Unlock()
			}
			break
		}

		err = ps.writeOpusFrame(packet)
		if err != nil {
			ps.Lock()
			ps.err = err
			ps.Unlock()
			break
		}
	}
}

func (ps *EncodingSession) writeOpusFrame(frame []byte) error {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.LittleEndian, int16(len(frame)))
	if err != nil {
		return err
	}

	_, err = buffer.Write(frame)
	if err != nil {
		return err
	}

	ps.data <- buffer.Bytes()
	ps.Lock()
	ps.lastFrame++
	ps.Unlock()
	return nil
}
