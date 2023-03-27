package runner

import (
	"io/ioutil"
	"os"
	"strconv"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"
	"github.com/sashabaranov/go-openai"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
)

func (r *Runner) CountTokens(cmd *console.Script) console.ExitCode {
	r.loadCfg(cmd, question.NewHelper(os.Stdin, output.NewCliOutput(true, nil)))
	var (
		filePath = cmd.Input.Argument("file")
		msgs     = []openai.ChatCompletionMessage{}
	)
	gauge, err := gpt.NewCostGauge(r.cfg.Model, entity.BpeEncoderFile, entity.BpeVocabFile)
	iferr.Exit(err)

	b, err := ioutil.ReadFile(filePath)
	iferr.Exit(err)

	msgs = append(msgs, openai.ChatCompletionMessage{
		Role:    "user",
		Name:    r.cfg.User,
		Content: string(b),
	})
	meta, err := gauge.CountTokensIn(msgs)
	iferr.Exit(err)
	cmd.PrintNote(strconv.Itoa(meta.NumTokens) + " tokens in file " + filePath)
	return console.ExitSuccess
}
