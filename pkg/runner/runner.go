package runner

import (
	"context"
	"errors"
	"fmt"
	"os"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/color"
	"github.com/DrSmithFr/go-console/formatter"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/DrSmithFr/go-console/question/answers"
	"github.com/go-dockly/utility/xerrors/iferr"
	"github.com/sashabaranov/go-openai"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
	"github.com/go-dockly/ggpt/pkg/util"
)

var banner = `
  ____        ____ ____ _____ 
 / ___| ___  / ___|  _ \_   _|
| |  _ / _ \| |  _| |_) || |  
| |_| | (_) | |_| |  __/ | |  
 \____|\___/ \____|_|    |_|  
							  
`

type IRunner interface {
	ShowChat(cmd *console.Script) console.ExitCode
	Completion(cmd *console.Script) console.ExitCode
	CountTokens(cmd *console.Script) console.ExitCode
	SetModel(cmd *console.Script) console.ExitCode
	UpdateSettings(cmd *console.Script) console.ExitCode
	Chat(cmd *console.Script) console.ExitCode
	IChat(cmd *console.Script) console.ExitCode
	Chats(cmd *console.Script) console.ExitCode
	DeleteChat(cmd *console.Script) console.ExitCode
	PrintPersonas(cmd *console.Script) console.ExitCode
	Recorder(cmd *console.Script) console.ExitCode
	Say(cmd *console.Script) console.ExitCode
	PlaySound(cmd *console.Script) console.ExitCode
	TranscribeFile(cmd *console.Script) console.ExitCode
	TranslateSQLQuery(cmd *console.Script) console.ExitCode
	TranslateToSQLQuery(cmd *console.Script) console.ExitCode
	ExecuteShellCommand(cmd *console.Script) console.ExitCode
	ExplainShellCommand(cmd *console.Script) console.ExitCode
	TranslateRegex(cmd *console.Script) console.ExitCode
	ExplainRegex(cmd *console.Script) console.ExitCode
	Cost(cmd *console.Script) console.ExitCode
	Encode(cmd *console.Script) console.ExitCode
}

type Runner struct {
	cfg     *entity.Config
	session *entity.Session
}

var (
	f1 = formatter.NewOutputFormatterStyle(color.Green, color.Default, []string{color.Underscore, color.Bold})
	f2 = formatter.NewOutputFormatterStyle(color.Default, color.Blue, nil)
	f3 = formatter.NewOutputFormatterStyle(color.Red, color.Default, nil)
)

func New() (IRunner, error) {
	return &Runner{}, nil
}

func (r *Runner) init(cmd *console.Script, qh *question.Helper) {

	if !util.FileExists(entity.SnowboyResFile) {
		_ = os.Mkdir(entity.GgptDir, os.ModePerm)

		answer := qh.Ask(
			question.NewComfirmation("Do you want to use snowboy for hotword detection?").
				SetDefaultAnswer(answers.Yes))
		if answer == answers.Yes {
			var repoURL = "https://github.com/Kitt-AI/snowboy/raw/master/"
			fmt.Println("snowboy hotword detector download started")

			iferr.Exit(util.DownloadFile(repoURL+"resources/common.res", entity.SnowboyDir, "common.res"))

			iferr.Exit(util.DownloadFile(repoURL+"resources/models/computer.umdl", entity.SnowboyDir, "computer.umdl"))

			iferr.Exit(util.DownloadFile(repoURL+"LICENSE", entity.SnowboyDir, "LICENSE"))

			fmt.Println("download finished")
		} else {
			iferr.Exit(errors.New("voice chat needs snowboy to be installed"))
		}
	}

	if !util.FileExists(entity.BpeEncoderFile) {
		var repoURL = "https://github.com/latitudegames/GPT-3-Encoder/raw/master/"
		fmt.Println("download bpe token encoder files")

		iferr.Exit(util.DownloadFile(repoURL+"encoder.json", entity.BpeDir, "encoder.json"))

		iferr.Exit(util.DownloadFile(repoURL+"vocab.bpe", entity.BpeDir, "vocab.bpe"))

		iferr.Exit(util.DownloadFile(repoURL+"LICENSE", entity.BpeDir, "LICENSE"))

		fmt.Println("download finished")
	}

}

func (r *Runner) Completion(cmd *console.Script) console.ExitCode {
	var (
		qh      = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
		persona = cmd.Input.Option("persona")
		prompt  = cmd.Input.Argument("query")
	)
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)
	if persona == "" {
		persona = "default"
	}
	_, err := r.complete(cmd, qh, prompt, persona, "false")
	iferr.Exit(err)

	r.saveSession(cmd, qh)
	return console.ExitSuccess
}

func (r *Runner) complete(cmd *console.Script, qh *question.Helper, prompt, persona, tts string) (answer string, err error) {
	var (
		ai      = gpt.NewTTSClient(r.cfg.APIKey, r.cfg.Model, tts)
		ctx     = context.Background()
		request = r.mapChatCompletionSettings()
	)

	request.Messages = append(r.loadPersona(cmd, persona), openai.ChatCompletionMessage{
		Role:    "user",
		Content: prompt,
	})

	request.Messages, err = ai.ChatCompletion(ctx, request)
	if err != nil {
		return answer, err
	}

	r.printTokenCount(cmd, request.Messages)
	answer = request.Messages[len(request.Messages)-1].Content

	return answer, nil
}
