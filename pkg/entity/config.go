package entity

import "os"

type Config struct {
	APIKey           string         `json:"api_key"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Model            string         `json:"model"`
	N                int            `json:"n,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	TTSEnabled       bool           `json:"tts_enabled"`
	User             string         `json:"user"`
	Echo             bool           `json:"echo,omitempty"`
	Suffix           string         `json:"suffix,omitempty"`
	BestOf           int            `json:"best_of,omitempty"`
	LogProbs         int            `json:"log_probs,omitempty"`
}

var (
	GgptDir         = os.Getenv("HOME") + "/.ggpt/data/"
	ConfigFile      = GgptDir + "config.json"
	SessionFile     = GgptDir + "session.json"
	SnowboyDir      = GgptDir + "snowboy/"
	SnowboyResFile  = SnowboyDir + "common.res"
	SnowboyUMDLFile = SnowboyDir + "computer.umdl"
	BpeDir          = GgptDir + "bpe/"
	BpeEncoderFile  = BpeDir + "encoder.json"
	BpeVocabFile    = BpeDir + "vocab.bpe"
)
