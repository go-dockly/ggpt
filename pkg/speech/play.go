package speech

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	errors "github.com/go-dockly/utility/xerrors"
)

func (s *Service) Play(fileName string) (err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", fileName)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return errors.Wrap(err, "failed to decode wav file")
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	speaker.Close()

	return nil
}
