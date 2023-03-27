package gpt

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/go-dockly/utility/xerrors/iferr"
	"github.com/jdkato/prose/v2"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

type IClient interface {
	TTS(msg string) error
	STT(fileName string) (string, error)
	ChatCompletion(ctx context.Context, request *openai.ChatCompletionRequest) ([]openai.ChatCompletionMessage, error)
}
type Client struct {
	ai    *openai.Client
	model string
}

func NewClient(apiKey, model string, tts ...string) IClient {
	return &Client{model: model, ai: openai.NewClient(apiKey)}
}

func NewTTSClient(apiKey, model string, tts ...string) IClient {
	c := &Client{model: model, ai: openai.NewClient(apiKey)}

	if len(tts) > 0 && tts[0] == "true" {
		err := speaker.Init(22050, 2205)
		if err != nil {
			fmt.Println("tts disabled: init speaker: ", err)
		} else {
			go c.ttsSaveLoop()
			go c.ttsPlayLoop()
		}
	}

	return c
}

func (c *Client) ChatCompletion(ctx context.Context, request *openai.ChatCompletionRequest) (messages []openai.ChatCompletionMessage, err error) {
	// spew.Dump(request)
	stream, err := c.ai.CreateChatCompletionStream(ctx, *request)
	if err != nil {
		return nil, errors.Wrap(err, "gpt completion")
	}
	defer stream.Close()

	var (
		counter   = 0
		answer    = ""
		segAnswer = ""
		replace   = ""
	)

	for {
		data, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println()
				break
			}
			return nil, fmt.Errorf("stream error: %v", err)
		} else {
			if len(data.Choices) > 0 {
				counter++
				answer += data.Choices[0].Delta.Content
				fmt.Print(aurora.Red(data.Choices[0].Delta.Content))
				segAnswer += data.Choices[0].Delta.Content
				replace = c.replaceText(segAnswer, true)
				segAnswer = strings.Replace(segAnswer, replace, "", 1)
			}
		}
	}
	if counter == 0 {
		return nil, fmt.Errorf("stream did not return any data")
	} else {
		c.replaceText(segAnswer, false)

		messages = append(request.Messages, openai.ChatCompletionMessage{
			Role:    "assistant",
			Content: answer,
		})
	}

	return messages, nil
}

var toSpeak = []string{}

func (c *Client) ttsSaveLoop() {
	var counter = 0
	for {
		if len(toSpeak) > 0 {
			var say = toSpeak[0]
			toSpeak = toSpeak[1:] // delete first element
			iferr.Warn(c.toFile(say, counter))
			counter++
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Client) ttsPlayLoop() {
	for {
		var dir = fmt.Sprintf("%s/.ggpt/data/microphone", os.Getenv("HOME"))
		files, err := ioutil.ReadDir(dir)
		iferr.Exit(err, "read dir")
		for _, file := range files {
			var filePath = fmt.Sprintf("%s/%s", dir, file.Name())
			err = c.play(filePath)
			if err != nil {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			iferr.Warn(os.Remove(filePath))
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Client) play(fileName string) (err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", fileName)
	}

	streamer, _, err := wav.Decode(f)
	if err != nil {
		return errors.Wrap(err, "failed to decode wav file")
	}
	defer streamer.Close()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done

	return nil
}

func (c *Client) replaceText(answer string, firstSentenceOnly bool) (replace string) {
	doc, err := prose.NewDocument(answer, prose.WithTokenization(false))
	iferr.Warn(err)

	sentences := doc.Sentences()
	if firstSentenceOnly {
		if len(sentences) > 1 {
			replace = sentences[0].Text
			toSpeak = append(toSpeak, replace)

			return replace
		}
		return ""
	} else {
		for _, v := range sentences {
			replace = v.Text
			toSpeak = append(toSpeak, replace)
		}
	}

	return replace
}
