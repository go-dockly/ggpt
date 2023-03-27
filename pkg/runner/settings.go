package runner

import (
	"fmt"
	"strconv"
	"strings"

	console "github.com/DrSmithFr/go-console"
	"github.com/go-dockly/utility/xerrors/iferr"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
	"github.com/go-dockly/ggpt/pkg/util"
)

func (r *Runner) loadProfile(profile string) {
	switch profile {
	case "creative":
		// combination of settings should result in more creative and less predictable text generation.
		r.cfg.Temperature = 0.7
		r.cfg.N = 1
		r.cfg.TopP = 0.9
		r.cfg.FrequencyPenalty = 1
		r.cfg.PresencePenalty = 0.6
		r.cfg.MaxTokens = 300
		r.cfg.Stop = []string{"\n\n"} // limit response to a single paragraph
	case "focused":
		r.cfg.Temperature = 0.5
		r.cfg.N = 1
		r.cfg.TopP = 0.5
		r.cfg.FrequencyPenalty = 1
		r.cfg.PresencePenalty = 0.6
		r.cfg.MaxTokens = 150               // responsive
		r.cfg.Stop = []string{"\n", "\n\n"} // limit response to a single paragraph
	case "default":
		r.cfg.Temperature = 1
		r.cfg.N = 1
		r.cfg.TopP = 1
		r.cfg.FrequencyPenalty = 0
		r.cfg.PresencePenalty = 0
		r.cfg.MaxTokens = 2048
		r.cfg.Stop = []string{}
	default:
		fmt.Printf("unknown profile: %s available options: creative/focused/default\n", profile)
	}
}

// todo load moderate / creative preset profiles
// presence_penalty=0.6 suggests a moderate level of penalty
// to discourage repetitions while still allowing some similarity in the generated text.
func (r *Runner) UpdateSettings(cmd *console.Script) console.ExitCode {
	iferr.Exit(util.LoadFile(entity.ConfigFile, &r.cfg), "failed to load config")

	var temperature = cmd.Input.Option("temperature")
	if temperature != "" {
		t, err := strconv.ParseFloat(temperature, 64)
		iferr.Exit(err, "failed to parse temperature")
		r.cfg.Temperature = float32(t)
	}

	var profile = cmd.Input.Option("profile")
	if profile != "" {
		r.loadProfile(profile)
	}

	var apiKey = cmd.Input.Option("api_key")
	if apiKey != "" {
		r.cfg.APIKey = apiKey
	}

	var user = cmd.Input.Option("user")
	if user != "" {
		r.cfg.User = user
	}

	var maxTokens = cmd.Input.Option("max_tokens")
	if maxTokens != "" {
		i, err := strconv.Atoi(maxTokens)
		iferr.Exit(err, "failed to parse max_tokens")
		r.cfg.MaxTokens = i
	}

	var topP = cmd.Input.Option("top_p")
	if topP != "" {
		t, err := strconv.ParseFloat(topP, 64)
		iferr.Exit(err, "failed to parse top_p")
		r.cfg.TopP = float32(t)
	}

	var n = cmd.Input.Option("n_completions")
	if n != "" {
		i, err := strconv.Atoi(n)
		iferr.Exit(err, "failed to parse n")
		r.cfg.N = i
	}

	var p = cmd.Input.Option("presence_penalty")
	if p != "" {
		t, err := strconv.ParseFloat(p, 64)
		iferr.Exit(err, "failed to parse presence_penalty")
		r.cfg.PresencePenalty = float32(t)
	}

	var f = cmd.Input.Option("frequency_penalty")
	if f != "" {
		t, err := strconv.ParseFloat(f, 64)
		iferr.Exit(err, "failed to parse frequency_penalty")
		r.cfg.FrequencyPenalty = float32(t)
	}

	// https://aidungeon.medium.com/controlling-gpt-3-with-logit-bias-55866d593292
	var l = cmd.Input.Option("logit_bias")
	if l != "" {
		gauge, err := gpt.NewCostGauge(r.cfg.Model, entity.BpeEncoderFile, entity.BpeVocabFile)
		iferr.Exit(err)
		var kvs = strings.Split(l, ",")
		for _, pairs := range kvs { // max length 300 biases
			// eg "Paris:-10"
			var pair = strings.Split(pairs, ":")
			i, err := strconv.Atoi(pair[1])
			iferr.Exit(err, "failed to parse logit_bias")
			_, err = strconv.Atoi(pair[0])
			if err != nil { // is it a plain text word
				tokens := gauge.Encode(" " + pair[0])
				token := strconv.Itoa(int(tokens[0]))
				if len(r.cfg.LogitBias) == 0 {
					r.cfg.LogitBias = make(map[string]int)
				}
				r.cfg.LogitBias[token] = i
			}
		}
	}

	var stop = cmd.Input.Option("stop")
	if stop != "" {
		r.cfg.Stop = strings.Split(stop, ",")
	}

	iferr.Exit(util.WriteFile(entity.ConfigFile, r.cfg), "failed to save config")

	return console.ExitSuccess
}
