package runner

import (
	"os"
	"strings"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"
)

func (r *Runner) ExplainRegex(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)
	var (
		command = cmd.Input.Argument("expression")
		prompt  = "explain the following regex: " + command
		persona = "default"
		tts     = "false"
	)

	_, err := r.complete(cmd, qh, prompt, persona, tts)
	iferr.Exit(err)

	r.saveSession(cmd, qh)

	return console.ExitSuccess
}

func (r *Runner) TranslateRegex(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)
	var (
		command = cmd.Input.Argument("expression")
		prompt  = r.buildRegexPrompt(command)
		persona = "default"
		tts     = "false"
	)

	_, err := r.complete(cmd, qh, prompt, persona, tts)
	iferr.Exit(err)

	r.saveSession(cmd, qh)

	return console.ExitSuccess
}

func (r *Runner) buildRegexPrompt(command string) string {
	var (
		explainText = "In addition, provide a description of how the regex works."
		formatText  = "Expression: <insert_expression_here>\nDescription: <insert_description_here>"
	)
	prompt := `You're a command line tool that generates regular expressions for the user.
Format: {format_text}
Instructions: {instruction}. {explain_text}`

	prompt = strings.Replace(prompt, "{format_text}", formatText, 1)
	prompt = strings.Replace(prompt, "{explain_text}", explainText, 1)
	prompt = strings.Replace(prompt, "{instruction}", command, 1)

	return prompt
}
