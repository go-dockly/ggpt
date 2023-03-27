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

func (r *Runner) PrintPersonas(cmd *console.Script) console.ExitCode {
	var (
		qh      = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
		persona = cmd.Input.Argument("persona")
		session = new(entity.Session)
	)
	iferr.Exit(util.LoadFile(entity.SessionFile, session), "failed to load personas")
	if persona != "" {
		r.printPersona(cmd, session, persona)
	} else {
		keys := make([]string, 0, len(session.Personas))
		for k := range session.Personas {
			keys = append(keys, k)
		}
		choice := qh.Ask(
			question.NewChoices("Existing personas", keys).
				SetMultiselect(false).
				SetMaxAttempts(1),
		)
		r.printPersona(cmd, session, choice)
	}

	return console.ExitSuccess
}

func (r *Runner) printPersona(cmd *console.Script, session *entity.Session, persona string) {
	profile, ok := session.Personas[persona]
	if ok {
		cmd.PrintSuccess(profile.Content)
		return
	}
	cmd.PrintError(persona + " not found")
}

// effectively resets the chat session
func (r *Runner) loadPersona(cmd *console.Script, persona string) (messages []openai.ChatCompletionMessage) {
	if persona != "default" {
		cmd.PrintText("Loading profile for " + persona)
	}
	profile, ok := r.session.Personas[persona]
	if ok {
		// todo intent
		// if persona == "intent" {
		// 	intents, err := modules.LoadPatterns()
		// 	iferr.Exit(err, "loading intent patterns")
		// 	r.session.Persona = persona
		// 	modules.Init() // initiallize pattern mapped modules
		// 	profile.Content = fmt.Sprintf(profile.Content, intents)
		// }
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "system",
			Content: profile.Content,
		})
	} else {
		cmd.PrintText(f3.Apply(persona + " profile not found - using default "))
		return r.loadPersona(cmd, "default")
	}
	return messages
}
