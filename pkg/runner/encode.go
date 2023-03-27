package runner

import (
	"fmt"
	"os"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
)

func (r *Runner) Encode(cmd *console.Script) console.ExitCode {
	var (
		qh   = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
		word = cmd.Input.Argument("word")
	)
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	gauge, err := gpt.NewCostGauge(r.cfg.Model, entity.BpeEncoderFile, entity.BpeVocabFile)
	iferr.Exit(err)
	tokens := gauge.Encode(" " + word)
	fmt.Println(tokens)

	return console.ExitSuccess
}
