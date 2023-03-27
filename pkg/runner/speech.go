package runner

import (
	"fmt"
	"os"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"
	"github.com/logrusorgru/aurora"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
	"github.com/go-dockly/ggpt/pkg/speech"
)

func (r *Runner) Recorder(cmd *console.Script) console.ExitCode {
	fmt.Println(aurora.Yellow(banner))
	var (
		qh           = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
		persona      = cmd.Input.Argument("persona")
		tts          = cmd.Input.Option("tts")
		printWelcome = true
	)
	r.init(cmd, qh)
	r.loadCfg(cmd, qh, printWelcome)
	r.loadSession(cmd, qh)
	var (
		ai = gpt.NewTTSClient(r.cfg.APIKey, r.cfg.Model, tts)
		s  = speech.New(ai)
	)

	r.session.CurrentChat = r.mapChatCompletionSettings()

	if persona != "" {
		r.session.CurrentChat.Messages = r.loadPersona(cmd, persona)
	}
	cmd.PrintText("Using gpt model " + r.cfg.Model)

	iferr.Exit(s.Listen(r.session, entity.SnowboyUMDLFile, entity.SnowboyResFile), "recorder listen")
	return console.ExitSuccess
}

func (r *Runner) Say(cmd *console.Script) console.ExitCode {
	r.loadCfg(cmd, question.NewHelper(os.Stdin, output.NewCliOutput(true, nil)))
	var (
		ai      = gpt.NewTTSClient(r.cfg.APIKey, r.cfg.Model)
		message = cmd.Input.Argument("message")
	)
	iferr.Exit(ai.TTS(message), message)
	return console.ExitSuccess
}

func (r *Runner) PlaySound(cmd *console.Script) console.ExitCode {
	var (
		fileName = cmd.Input.Argument("file")
		s        = speech.New(nil)
	)
	iferr.Exit(s.Play(fileName), "play file")
	return console.ExitSuccess
}

func (r *Runner) TranscribeFile(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)
	var (
		fileName = cmd.Input.Argument("file")
		ai       = gpt.NewClient(r.cfg.APIKey, r.cfg.Model)
	)
	text, err := ai.STT(fileName)
	iferr.Exit(err, "stt")
	cmd.PrintSuccess(text)
	r.saveSession(cmd, qh)
	return console.ExitSuccess
}
