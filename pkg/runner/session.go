package runner

import (
	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"
	"github.com/sashabaranov/go-openai"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
	"github.com/go-dockly/ggpt/pkg/util"
)

func (r *Runner) loadSession(cmd *console.Script, qh *question.Helper) {
	r.session = new(entity.Session)
	if !util.FileExists(entity.SessionFile) {
		r.session.Persona = "default"
		r.session.Cost = make(map[string]*gpt.CostGaugeMeta)
		r.session.Personas = make(map[string]*entity.Persona)
		r.session.Personas["default"] = &entity.Persona{
			Content: "When formulating an answer do not abbreviate words.",
		}
		r.session.CurrentChat = new(openai.ChatCompletionRequest)
		iferr.Exit(util.WriteFile(entity.SessionFile, r.session))
		cmd.PrintText("session saved")
	}
	iferr.Exit(util.LoadFile(entity.SessionFile, r.session))
}

func (r *Runner) saveSession(cmd *console.Script, qh *question.Helper) {
	// var today = time.Now().Format("06-Jan-02")
	// r.session.Cost[today] = meta
	r.session.CurrentChat = new(openai.ChatCompletionRequest)
	// answer := qh.Ask(
	// 	question.NewComfirmation("Do you want to save the session?").
	// 		SetDefaultAnswer(answers.Yes).
	// 		SetMaxAttempts(2),
	// )
	// if answer == answers.Yes {
	// 	r.session.Chats[] = session.CurrentChat
	// 	cmd.PrintText("session saved")
	// } else {
	// 	cmd.PrintText("... ok :)")
	// }
	iferr.Exit(util.WriteFile(entity.SessionFile, r.session))
}
