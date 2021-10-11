package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/jonas747/ogg"
)

// https://github.com/jonas747/dca

type EncodingSession struct {
	sync.Mutex

	path   string    // Encoding from Path
	reader io.Reader // Encoding from Reader

	running   bool
	process   *os.Process
	lastFrame int

	err    error
	stderr bytes.Buffer

	data   chan []byte
	buffer bytes.Buffer
}

func NewEncodingFromPath(path string) *EncodingSession {
	session := &EncodingSession{path: path, data: make(chan []byte)}
	go session.start()
	return session
}

func NewEncodingFromReader(reader io.Reader) *EncodingSession {
	session := &EncodingSession{reader: reader, data: make(chan []byte)}
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
	source := utils.Is(s.reader != nil, "pipe:0", s.path)

	arguments := []string{
		"-i", source, "-loglevel", "fatal",
		"-map", "0:a", "-acodec", "libopus",
		"-f", "ogg", "-frame_duration", "20", "pipe:1"}

	if source != "pipe:0" {
		arguments = append([]string{
			"-reconnect", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "2"}, arguments...)
	}
	println(source)

	cmd := exec.Command("ffmpeg", arguments...)
	if s.reader != nil {
		cmd.Stdin = s.reader
	}

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
	s.readStdout(stdout)

	defer func() {
		if s.err == nil {
			stderr := s.stderr.String()
			if stderr != "" && stderr != "<nil>" {
				s.err = errors.New(utils.Fmt("ffmpeg: %s", stderr))
			}
		}
	}()
	cmd.Wait()
}

func (s *EncodingSession) Stop() error {
	s.Lock()
	defer s.Unlock()

	if !s.running || s.process == nil {
		return errors.New("ffmpeg not running")
	}

	err := s.process.Kill()
	return err
}

func (s *EncodingSession) StopClean() error {
	err := s.Stop()

	if err == nil {
		for range s.data {
		}
	}
	return err
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
