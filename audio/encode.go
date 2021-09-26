package audio

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/jonas747/ogg"
)

// https://github.com/jonas747/dca

type ProcessingSession struct {
	sync.Mutex
	source    string
	running   bool
	process   *os.Process
	lastFrame int
	isOpus    bool

	err error

	data   chan []byte
	buffer bytes.Buffer
}

func NewProcessingSession(source string, isOpus bool) *ProcessingSession {
	session := &ProcessingSession{
		source: source,
		isOpus: isOpus,
		data:   make(chan []byte),
	}
	go session.start()
	return session
}

func (ps *ProcessingSession) start() {
	defer func() {
		ps.Lock()
		ps.running = false
		ps.Unlock()
	}()

	ps.Lock()
	ps.running = true

	arguments := []string{
		"-reconnect", "1", "-reconnect_at_eof", "1",
		"-reconnect_streamed", "1", "-reconnect_delay_max", "2",
		"-i", ps.source, "-loglevel", "fatal",
		"-map", "0:a", "-acodec", utils.Is(ps.isOpus, "copy", "libopus"),
		"-f", "ogg", "-ar", "48000",
		"-ac", "2", "-b:a", "96000",
		"-application", "audio", "-frame_duration", "20", "pipe:1"}

	process := exec.Command("ffmpeg", arguments...)
	stdout, err := process.StdoutPipe()
	if err != nil {
		ps.err = err
		ps.Unlock()
		close(ps.data)
		return
	}

	stderr, err := process.StderrPipe()
	if err != nil {
		ps.err = err
		ps.Unlock()
		close(ps.data)
		return
	}

	err = process.Start()
	if err != nil {
		ps.err = err
		ps.Unlock()
		close(ps.data)
		return
	}
	ps.process = process.Process
	ps.Unlock()

	var wg sync.WaitGroup
	wg.Add(1)
	go ps.readStderr(stderr, &wg)

	defer close(ps.data)
	ps.readStdout(stdout)

	process.Wait()
}

func (ps *ProcessingSession) Stop() error {
	ps.Lock()
	defer ps.Unlock()
	if !ps.running || ps.process == nil {
		return errors.New("not running")
	}

	err := ps.process.Kill()
	return err
}

func (ps *ProcessingSession) StopClean() {
	ps.Stop()
	for range ps.data {
	}
}

func (ps *ProcessingSession) FrameDuration() time.Duration {
	return time.Duration(20) * time.Millisecond
}

func (ps *ProcessingSession) OpusFrame() ([]byte, error) {
	frame := <-ps.data
	if frame == nil {
		return nil, io.EOF
	}

	if len(frame) < 2 {
		return nil, errors.New("bad frame")
	}
	return frame[2:], nil
}

func (ps *ProcessingSession) ReadFrame() ([]byte, error) {
	frame := <-ps.data
	if frame == nil {
		return nil, io.EOF
	}

	return frame, nil

}

func (ps *ProcessingSession) Read(p []byte) (int, error) {
	if ps.buffer.Len() >= len(p) {
		return ps.buffer.Read(p)
	}

	for ps.buffer.Len() < len(p) {
		frame, err := ps.ReadFrame()
		if err != nil {
			break
		}
		ps.buffer.Write(frame)
	}

	return ps.buffer.Read(p)
}

func (ps *ProcessingSession) readStderr(std io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()

	reader := bufio.NewReader(std)
	var buffer bytes.Buffer
	for {
		r, _, err := reader.ReadRune()

		if err != nil {
			if err != io.EOF {
				ps.Lock()
				ps.err = err
				ps.Unlock()
			}
			break
		}

		switch r {
		case '\n':
			str := strings.TrimSpace(strings.ReplaceAll(buffer.String(), ps.source, "source"))
			if str != "source: I/O error" {
				ps.Lock()
				println(str)
				ps.err = errors.New(str)
				ps.Unlock()
			}
			buffer.Reset()
		default:
			buffer.WriteRune(r)
		}
	}

}

func (ps *ProcessingSession) readStdout(std io.ReadCloser) {
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

func (ps *ProcessingSession) writeOpusFrame(frame []byte) error {
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
