package runner

import (
	"os"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"
	"github.com/sashabaranov/go-openai"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/util"
)

var supportedModels = []string{openai.GPT4, openai.GPT3Dot5Turbo, openai.GPT3Davinci}

func (r *Runner) SetModel(cmd *console.Script) console.ExitCode {
	iferr.Exit(util.LoadFile(entity.ConfigFile, &r.cfg), "failed to load config")
	model := cmd.Input.Argument("model")
	if model == "" {
		cmd.PrintText("current gpt model: " + r.getModel())
	}
	if util.Contains(supportedModels, model) {
		r.cfg.Model = model
	}
	iferr.Exit(util.WriteFile(entity.ConfigFile, r.cfg), "failed to save config")
	return console.ExitSuccess
}

func (r *Runner) getModel() string {
	if r.cfg.Model == "" {
		var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
		r.cfg.Model = qh.Ask(
			question.NewChoices("Which model to load?", supportedModels).
				SetMultiselect(false).
				SetMaxAttempts(1),
		)
		iferr.Exit(util.WriteFile(entity.ConfigFile, r.cfg), "failed to save config")
	}

	return r.cfg.Model
}
