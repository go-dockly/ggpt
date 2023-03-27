package speech

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/brentnd/go-snowboy"
	"github.com/go-dockly/utility/xerrors/iferr"
	"github.com/gordonklaus/portaudio"
	"github.com/logrusorgru/aurora"
	"github.com/sashabaranov/go-openai"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
	"github.com/go-dockly/ggpt/pkg/util"
)

func (s *Service) Listen(session *entity.Session, hotwordFile, detectorFile string) error {
	var (
		ctx             = context.Background()
		sig             = make(chan bool, 1)
		hotWordDetected = false
		silenceCounter  = 0
		inputChannels   = 1
		outputChannels  = 0
		sampleRate      = 16000
		framesPerBuffer = 1024
		frames          = make([]int16, framesPerBuffer)
		countTokens     = true
		currentMonth    = time.Now().Format("06-Jan")
	)
	today, ok := session.Cost[currentMonth]
	if !ok {
		today = new(gpt.CostGaugeMeta)
	}

	gauge, err := gpt.NewCostGauge(session.CurrentChat.Model, entity.BpeEncoderFile, entity.BpeVocabFile)
	if err != nil {
		fmt.Println("token counter disabled: " + err.Error())
		countTokens = false
	}

	err = portaudio.Initialize()
	if err != nil {
		return err
	}
	defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), framesPerBuffer, frames)
	if err != nil {
		return err
	}
	defer stream.Close()
	var sounder = NewSound(stream, frames)

	d := snowboy.NewDetector(detectorFile)
	defer d.Close()

	d.HandleFunc(snowboy.NewHotword(hotwordFile, 0.5), func(res string) {
		sounder.StartRecording()
		fmt.Println(aurora.BrightGreen("Hotword detected"))
		hotWordDetected = true
		// iferr.Warn(s.Play("data/wav/select.wav"))
		go func() {
			for {
				select {
				case <-sig:
					fmt.Println(aurora.Yellow("stop recording"))
					seconds, err := sounder.StopRecording()
					iferr.Exit(err)

					// iferr.Warn(s.Play("data/wav/select.wav"))
					hotWordDetected = false

					prompt, err := s.ai.STT(tmpWavFile)
					iferr.Exit(err)
					fmt.Println(prompt)

					session.CurrentChat.Messages = append(session.CurrentChat.Messages, openai.ChatCompletionMessage{
						Role:    "user",
						Content: prompt,
					})

					if countTokens {
						gauge.AddWhisperSeconds(seconds)
						delta, err := gauge.CountTokensIn(session.CurrentChat.Messages)
						iferr.Warn(err)
						session.Cost[currentMonth] = gauge.UpdateCost(today, delta)
						gauge.Reset()
						iferr.Exit(util.WriteFile(entity.SessionFile, session))
					}

					session.CurrentChat.Messages, err = s.ai.ChatCompletion(ctx, session.CurrentChat)
					iferr.Warn(err)

					// lenMsg := len(session.CurrentChat.Messages)
					// var answer = session.CurrentChat.Messages[lenMsg-1].Content
					// if session.Persona == "intent" {
					// 	answer = modules.ReplaceContent(answer)
					// 	fmt.Println(answer)
					// 	session.CurrentChat.Messages = s.reset(session)
					// }

					sounder.StartRecording()

					return
				}
			}
		}()
	})

	d.HandleSilenceFunc(500*time.Millisecond, func(string) {
		// fmt.Println("silence detected")
		if hotWordDetected {
			silenceCounter++
			if silenceCounter == 3 {
				hotWordDetected = false
				silenceCounter = 0
				sig <- true
			}
		}
	})

	sr, _, bd := d.AudioFormat()
	re := regexp.MustCompile(`/(\w+)\.umdl`)
	hotword := re.FindStringSubmatch(hotwordFile)[1]
	fmt.Printf("microphone streaming: sample_rate=%d, bit_depth=%d, hotword=%s, silence_counter=1.5sec\n", sr, bd, hotword)

	err = stream.Start()
	if err != nil {
		return err
	}

	err = d.ReadAndDetect(sounder)
	if err != nil {
		return err
	}

	return nil
}

// func (s *Service) reset(session *entity.Session) (messages []openai.ChatCompletionMessage) {
// 	var profile = session.Personas["intent"]
// 	intents, err := modules.LoadPatterns()
// 	iferr.Exit(err, "loading patterns")
// 	profile.Content = fmt.Sprintf(profile.Content, intents)
// 	return append(messages, openai.ChatCompletionMessage{
// 		Role:    "system",
// 		Content: profile.Content,
// 	})
// }
