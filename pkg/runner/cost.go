package runner

import (
	"fmt"
	"os"
	"time"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
)

func (r *Runner) Cost(cmd *console.Script) console.ExitCode {
	r.loadSession(cmd, question.NewHelper(os.Stdin, output.NewCliOutput(true, nil)))
	var month = cmd.Input.Argument("month")
	if month == "" {
		month = time.Now().Format("06-Jan")
	}
	v, ok := r.session.Cost[month]
	if ok {
		fmt.Printf("cost tokens: $%.3f\n", v.TotalCostGeneratedText)
		fmt.Printf("total tokens processed %d\n", v.TotalTokensGenerated)
		fmt.Printf("cost whisper $%.3f\n", v.TotalCostTranscribed)
		fmt.Printf("total seconds transcribed %.2f\n", v.TotalSecondsTranscribed)
	} else {
		return console.ExitError
	}
	return console.ExitSuccess
}
