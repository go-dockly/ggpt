package runner

import (
	"os"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"
	"github.com/sashabaranov/go-openai"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/util"
)

func (r *Runner) loadSchema(filePath string) string {
	b, err := os.ReadFile(filePath) // just pass the file name
	iferr.Exit(err, "failed to load schema")
	return string(b)
}

func (r *Runner) loadCfg(cmd *console.Script, qh *question.Helper, printWelcome ...bool) {
	var cfg = new(entity.Config)
	_ = util.LoadFile(entity.ConfigFile, cfg)
	r.cfg = cfg
	if r.cfg.User == "" {
		userName := qh.Ask(
			question.NewQuestion("What is your nickname?").
				SetDefaultAnswer("friend"),
		)
		cmd.PrintText("Hello " + userName)
		r.cfg.User = userName
		r.askModelSettings(cmd, qh)
		r.saveCfg(cmd)
	}
	if len(printWelcome) > 0 {
		cmd.PrintText(f1.Apply("Welcome back " + r.cfg.User))
	}
}

func (r *Runner) askModelSettings(cmd *console.Script, qh *question.Helper) {
	r.cfg.Model = qh.Ask(
		question.NewChoices("Which model to load?", supportedModels).
			SetMultiselect(false).
			SetMaxAttempts(2),
	)
	r.cfg.APIKey = qh.Ask(
		question.NewQuestion("What is your api key?"),
	)
}

func (r *Runner) saveCfg(cmd *console.Script) {
	iferr.Exit(util.WriteFile(entity.ConfigFile, r.cfg))
	cmd.PrintText("config saved")
}

func (r *Runner) mapChatCompletionSettings() *openai.ChatCompletionRequest {
	return &openai.ChatCompletionRequest{
		User:             r.cfg.User,
		Stream:           true,
		Model:            r.cfg.Model,
		Temperature:      r.cfg.Temperature,
		MaxTokens:        r.cfg.MaxTokens,
		LogitBias:        r.cfg.LogitBias,
		PresencePenalty:  r.cfg.PresencePenalty,
		FrequencyPenalty: r.cfg.FrequencyPenalty,
		Stop:             r.cfg.Stop,
		TopP:             r.cfg.TopP,
		N:                r.cfg.N,
	}
}

func (r *Runner) mapCompletionSettings() *openai.CompletionRequest {
	return &openai.CompletionRequest{
		User:             r.cfg.User,
		Stream:           true,
		Model:            r.cfg.Model, // todo validate that model is not higher than davinci for completion
		Temperature:      r.cfg.Temperature,
		MaxTokens:        r.cfg.MaxTokens,
		LogitBias:        r.cfg.LogitBias,
		PresencePenalty:  r.cfg.PresencePenalty,
		FrequencyPenalty: r.cfg.FrequencyPenalty,
		Stop:             r.cfg.Stop,
		TopP:             r.cfg.TopP,
		N:                r.cfg.N,
		Echo:             r.cfg.Echo,
		Suffix:           r.cfg.Suffix,
		BestOf:           r.cfg.BestOf,
		LogProbs:         r.cfg.LogProbs,
	}
}
