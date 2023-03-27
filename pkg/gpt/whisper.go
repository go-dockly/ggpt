package gpt

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

func (c *Client) STT(fileName string) (string, error) {
	var (
		ctx     = context.Background()
		request = openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: fileName,
		}
	)

	resp, err := c.ai.CreateTranscription(ctx, request)
	if err != nil {
		return "", err
	}

	return resp.Text, nil
}

// func (c *Client) play(msg string) (err error) {
// 	var (
// 		osName   = runtime.GOOS
// 		osPlayer = "aplay -"
// 	)
// 	switch osName {
// 	case "darwin":
// 		osPlayer = "ffplay -autoexit -nodisp -"
// 	case "windows":
// 		return errors.New("tts: windows not supported yet")
// 	}
// 	msg = normalizePunctuation(msg)
// 	cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -G --output - --data-urlencode 'text=%s' 'http://localhost:5002/api/tts' | %s", msg, osPlayer))
// 	err = cmd.Run()
// 	if err != nil {
// 		return errors.Wrap(err, "play")
// 	}

// 	return nil
// }

func (c *Client) toFile(msg string, counter int) error {
	if len(msg) > 250 {
		return fmt.Errorf("tts: message is too long")
	}

	url := "http://localhost:5002/api/tts?text=" + url.QueryEscape(normalizePunctuation(msg))
	filePath := fmt.Sprintf("%s/.ggpt/data/microphone/record%d.wav", os.Getenv("HOME"), counter)
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) TTS(sentence string) error {
	if len(sentence) > 250 {
		return fmt.Errorf("tts: message is too long")
	}

	var path = "http://localhost:5002/api/tts"

	r, err := http.Get(path + "?text=" + url.QueryEscape(normalizePunctuation(sentence)))
	if err != nil {
		return errors.Wrapf(err, "could not connect to mozilla/tacotron container on %s", path)
	}
	defer r.Body.Close()

	streamer, format, err := wav.Decode(r.Body)
	if err != nil {
		return errors.Wrap(err, "cannot decode wav file")
	}
	_ = format

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	streamer.Close()
	speaker.Close()
	return nil
}

func normalizePunctuation(msg string) string {
	var r = strings.NewReplacer("'", "", "%s", ".", "!?", "!", "!!", "!", "!!!", "!", "?!", "?", "??", "?", "???", "?", "..", ".", "...", ".")
	return r.Replace(msg)
}
