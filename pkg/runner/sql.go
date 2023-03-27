package runner

import (
	"fmt"
	"os"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/go-dockly/utility/xerrors/iferr"
)

func (r *Runner) TranslateToSQLQuery(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)
	var (
		query        = cmd.Input.Argument("text")
		schema       = cmd.Input.Argument("schema")
		schemaPrompt = ``
		persona      = "default"
		tts          = "false"
	)

	if schema != "" {
		schemaPrompt = fmt.Sprintf(`Use this table schema:\n\n%s\n\n`, r.loadSchema(schema))
	}

	var prompt = fmt.Sprintf(`Translate this natural language query into SQL:\n\n"%s"\n\n%sSQL Query:`, query, schemaPrompt)
	_, err := r.complete(cmd, qh, prompt, persona, tts)
	iferr.Exit(err)

	r.saveSession(cmd, qh)

	return console.ExitSuccess
}

func (r *Runner) TranslateSQLQuery(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)

	var (
		query   = cmd.Input.Argument("query")
		persona = "default"
		tts     = "false"
		prompt  = fmt.Sprintf(`Translate this SQL query into natural language:\n\n"%s"\n\nNatural language query:`, query)
	)

	_, err := r.complete(cmd, qh, prompt, persona, tts)
	iferr.Exit(err)

	r.saveSession(cmd, qh)

	return console.ExitSuccess
}
