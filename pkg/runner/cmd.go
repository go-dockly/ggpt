package runner

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"
)

func (r *Runner) ExplainShellCommand(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)
	var (
		command = cmd.Input.Argument("command")
		explain = true
		prompt  = r.buildCommandPrompt(command, explain)
		persona = "default"
		tts     = "false"
	)

	_, err := r.complete(cmd, qh, prompt, persona, tts)
	iferr.Exit(err)

	r.saveSession(cmd, qh)

	return console.ExitSuccess
}

func (r *Runner) ExecuteShellCommand(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)

	var (
		command = cmd.Input.Argument("command")
		explain = false
		prompt  = r.buildCommandPrompt(command, explain)
		persona = "default"
		tts     = "false"
	)

	answer, err := r.complete(cmd, qh, prompt, persona, tts)
	iferr.Exit(err)

	prompt = strings.Replace(answer, "Command: ", "", 1)
	prompt = strings.Replace(prompt, "\n", "", -1)

	out, err := exec.Command(os.Getenv("SHELL"), "-c", prompt).Output()
	iferr.Exit(err, prompt)

	fmt.Println(string(out))

	r.saveSession(cmd, qh)

	return console.ExitSuccess
}

func (r *Runner) buildCommandPrompt(command string, explain bool) string {
	var (
		explainText = "DO NOT EXPLAIN THE COMMAND"
		formatText  = "Command: <insert_command_here>"
		osName      = runtime.GOOS
	)
	if explain {
		explainText = "In addition, provide a detailed description of how the provided command works."
		formatText = "Command: <insert_command_here>\n Description: <insert_description_here>"
	}
	prompt := `You're a command line tool that generates CLI commands for the user.
Format: {format_text}
Instructions: Write a CLI command that does the following: {instruction}. It must work on {os} using {shell}. {explain_text}`

	prompt = strings.Replace(prompt, "{format_text}", formatText, 1)
	prompt = strings.Replace(prompt, "{explain_text}", explainText, 1)
	prompt = strings.Replace(prompt, "{os}", osName, 1)
	prompt = strings.Replace(prompt, "{shell}", os.Getenv("SHELL"), 1)
	prompt = strings.Replace(prompt, "{instruction}", command, 1)

	return prompt
}
