package speech

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	errors "github.com/go-dockly/utility/xerrors"
	"github.com/gordonklaus/portaudio"
)

type Sound struct {
	stream  *portaudio.Stream
	data    []int16
	file    *os.File
	encoder *wav.Encoder
	record  bool
}

const tmpWavFile = "data/record.wav"

func NewSound(stream *portaudio.Stream, data []int16) *Sound {
	file, _ := os.Create(tmpWavFile)
	wavEncoder := wav.NewEncoder(file, 16000, 16, 1, 1)

	return &Sound{
		file:    file,
		stream:  stream,
		data:    data,
		encoder: wavEncoder,
	}
}

func (s *Sound) StartRecording() {
	file, _ := os.Create(tmpWavFile)
	// todo don't write file to disk (openai marks the file invalid if not flushed to disk tho)
	s.encoder = wav.NewEncoder(file, 16000, 16, 1, 1)
	s.record = true
}

func (s *Sound) StopRecording() (seconds float64, err error) {
	err = s.encoder.Close()
	if err != nil {
		fmt.Println("Error closing WAV encoder:", err)
	}
	s.record = false

	return s.WavLength()
}

func (s *Sound) Read(p []byte) (int, error) {
	s.stream.Read()

	buf := &bytes.Buffer{}
	for _, v := range s.data {
		binary.Write(buf, binary.LittleEndian, v)
	}
	copy(p, buf.Bytes())

	if s.record {
		buffer := &audio.IntBuffer{
			Format: &audio.Format{
				NumChannels: 1,
				SampleRate:  16000,
			},
			Data:           int16ToInt(s.data),
			SourceBitDepth: 16,
		}
		err := s.encoder.Write(buffer)
		if err != nil {
			fmt.Println("Error writing to WAV file:", err)
			return 0, err
		}
	}

	return len(p), nil
}

func int16ToInt(slice []int16) []int {
	result := make([]int, len(slice))
	for i, value := range slice {
		result[i] = int(value)
	}
	return result
}

func (s *Sound) WavLength() (seconds float64, err error) {
	var (
		sampleRate    = 16000.0
		channels      = 1.0
		bitsPerSample = 16.0
	)
	file, err := s.file.Stat()
	if err != nil {
		return seconds, errors.Wrapf(err, "file info failed")
	}
	seconds = float64(file.Size()) / (sampleRate * channels * bitsPerSample / 8)
	var cost = seconds * 0.0001 // Whisper	$0.006 / minute $0.0001 / second (rounded to the nearest second)
	fmt.Printf("The recording is %.2f seconds long. whisper cost of %.3f$\n", seconds, cost)

	return seconds, nil
}
