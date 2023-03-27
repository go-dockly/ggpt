package gpt

import (
	"fmt"
	"os"

	"github.com/cohere-ai/tokenizer"
	"github.com/sashabaranov/go-openai"
)

type ICostGauge interface {
	AddWhisperSeconds(seconds float64)
	GetMaxTokenFor(model string) int
	CountTokensIn(messages []openai.ChatCompletionMessage) (meta *CostGaugeMeta, err error)
	UpdateCost(today *CostGaugeMeta, current *CostGaugeMeta) (total *CostGaugeMeta)
	Reset()
	Encode(word string) (tokens []int64)
}

func NewCostGauge(model, encoderFile, vocabFile string) (ICostGauge, error) {
	encoderReader, err := os.Open(encoderFile)
	if err != nil {
		return nil, err
	}
	vocabReader, err := os.Open(vocabFile)
	if err != nil {
		return nil, err
	}
	encoder, err := tokenizer.NewFromReaders(encoderReader, vocabReader)
	if err != nil {
		return nil, err
	}

	return &CostGauge{
		encoder: encoder,
		model:   model,
		delta:   new(CostGaugeMeta),
	}, nil
}

type CostGauge struct {
	model   string
	encoder *tokenizer.Encoder
	delta   *CostGaugeMeta
}

// todo Image Resolution Price gauge
// 1024×1024 $0.020 / image
// 512×512	$0.018 / image
// 256×256	$0.016 / image

type CostGaugeMeta struct {
	NumTokens               int
	Percentage              float64
	TotalSecondsTranscribed float64 `json:"whisper_seconds_transcribed"`
	TotalCostTranscribed    float64 `json:"whisper_cost_transcribed"`
	TotalTokensGenerated    int     `json:"gpt_total_tokens"`
	TotalCostGeneratedText  float64 `json:"gpt_text_cost"`
}

type GPTVersionTokenMap map[string]int

var gptVersionTokenMap = GPTVersionTokenMap{
	openai.GPT3Dot5Turbo: 4096,
	openai.GPT4:          8192,
	openai.GPT432K:       32768,
}

func (g *CostGauge) Encode(word string) []int64 {
	t, _ := g.encoder.Encode(word)
	return t
}

func (g *CostGauge) Reset() {
	g.delta = new(CostGaugeMeta)
}

func (g *CostGauge) GetMaxTokenFor(model string) int {
	return gptVersionTokenMap[model]
}

func (g *CostGauge) CountTokensIn(messages []openai.ChatCompletionMessage) (*CostGaugeMeta, error) {
	numTokens := 0
	for _, message := range messages {
		// underlying format consumed by ChatGPT models
		// https://github.com/openai/openai-python/blob/main/chatml.md
		numTokens += 4 // every message follows <im_start>{role/name}\n{content}<im_end>\n
		_, tokenStrings := g.encoder.Encode(message.Content)
		numTokens += len(tokenStrings)
		if message.Name != "" { // if there's a name, the role is omitted
			numTokens += -1 // otherwise role is always 1 token
		}
	}
	numTokens += 2 // every reply is primed with <im_start>assistant
	// gpt-3.5-turbo has 4,096 tokens (which is ~3 pages of single-lined English text).
	maxToken, ok := gptVersionTokenMap[g.model]
	if !ok {
		return nil, fmt.Errorf("unknown model: %s", g.model)
	}
	g.delta.NumTokens += numTokens
	g.delta.TotalTokensGenerated += numTokens
	g.delta.Percentage += (float64(numTokens) / float64(maxToken)) * 100
	g.delta.TotalCostGeneratedText += (float64(numTokens) * 0.000002)
	return g.delta, nil
}

func (g *CostGauge) AddWhisperSeconds(seconds float64) {
	g.delta.TotalSecondsTranscribed += seconds
	g.delta.TotalCostTranscribed += seconds * 0.0001
}

func (g *CostGauge) UpdateCost(month, today *CostGaugeMeta) (cost *CostGaugeMeta) {
	return &CostGaugeMeta{
		TotalSecondsTranscribed: month.TotalSecondsTranscribed + today.TotalSecondsTranscribed,
		TotalCostTranscribed:    month.TotalCostTranscribed + today.TotalCostTranscribed,
		TotalTokensGenerated:    month.TotalTokensGenerated + today.TotalTokensGenerated,
		TotalCostGeneratedText:  month.TotalCostGeneratedText + today.TotalCostGeneratedText,
	}
}
