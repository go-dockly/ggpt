package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/output"
	"github.com/DrSmithFr/go-console/question"
	"github.com/logrusorgru/aurora"
	"github.com/sashabaranov/go-openai"

	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
)

func (r *Runner) ShowChat(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.loadSession(cmd, qh)
	var name = cmd.Input.Argument("name")
	chat, ok := r.session.Chats[name]
	if ok {
		printMessages(chat, cmd)
	} else {
		cmd.PrintText(name + " not a known chat.")
		fmt.Printf("known chats are: %v\n", reflect.ValueOf(r.session.Chats).MapKeys())
	}
	return console.ExitSuccess
}

func (r *Runner) DeleteChat(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.loadSession(cmd, qh)
	var name = cmd.Input.Argument("name")
	_, ok := r.session.Chats[name]
	if ok {
		cmd.PrintText("delete chat " + name)
		delete(r.session.Chats, name)
	} else {
		cmd.PrintText("chat " + name + " not found")
	}
	r.saveSession(cmd, qh)
	return console.ExitSuccess

}

func (r *Runner) Chat(cmd *console.Script) console.ExitCode {
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)
	r.chat(cmd, qh)
	r.saveSession(cmd, qh)
	return console.ExitSuccess
}

func (r *Runner) IChat(cmd *console.Script) console.ExitCode {
	fmt.Println(aurora.Yellow(banner))
	var qh = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
	r.init(cmd, qh)
	r.loadCfg(cmd, qh)
	r.loadSession(cmd, qh)
	var persona = cmd.Input.Argument("persona")
	r.session.CurrentChat = r.mapChatCompletionSettings()
	if persona != "" {
		r.session.CurrentChat.Messages = r.loadPersona(cmd, persona)
	}
	r.loop(cmd, qh)
	r.saveSession(cmd, qh)
	return console.ExitSuccess
}

// TODO summarize conversation when approaching max tokens?
func (r *Runner) Chats(cmd *console.Script) console.ExitCode {
	var (
		qh           = question.NewHelper(os.Stdin, output.NewCliOutput(true, nil))
		printWelcome = true
	)
	r.init(cmd, qh)
	r.loadCfg(cmd, qh, printWelcome)
	r.loadSession(cmd, qh)
	var previous = cmd.Input.Argument("name")
	if previous != "" {
		prevChat, ok := r.session.Chats[previous]
		if ok {
			r.session.CurrentChat = prevChat
		} else {
			cmd.PrintNote(fmt.Sprintf("%s does not exist. Listing chats instead", previous))
			r.loadPrevChat(cmd, qh)
		}
	} else {
		r.loadPrevChat(cmd, qh)
	}
	r.loop(cmd, qh)
	r.saveSession(cmd, qh)
	return console.ExitSuccess
}

func (r *Runner) chat(cmd *console.Script, qh *question.Helper) (err error) {
	var (
		name   = cmd.Input.Argument("name")
		prompt = cmd.Input.Option("user")
		system = cmd.Input.Option("system")
		tts    = cmd.Input.Option("tts")
		ai     = gpt.NewTTSClient(r.cfg.APIKey, r.cfg.Model, tts)
		ctx    = context.Background()
	)

	_, ok := r.session.Chats[name]
	if !ok {
		r.session.Chats[name] = r.mapChatCompletionSettings()
		r.session.Chats[name].Messages = r.loadPersona(cmd, "default")
	}

	if system != "" {
		r.session.Chats[name].Messages = append(r.session.Chats[name].Messages, openai.ChatCompletionMessage{
			Role:    "system",
			Content: prompt,
		})
		cmd.PrintText("system prompt updated")
	}
	if prompt != "" {
		r.session.Chats[name].Messages = append(r.session.Chats[name].Messages, openai.ChatCompletionMessage{
			Role:    "user",
			Content: prompt,
		})
		r.session.Chats[name].Messages, err = ai.ChatCompletion(ctx, r.session.Chats[name])
		if err != nil {
			return err
		}
		r.printTokenCount(cmd, r.session.Chats[name].Messages)
	}

	return nil
}

func (r *Runner) loop(cmd *console.Script, qh *question.Helper) (err error) {
	cmd.PrintText("Using gpt model " + r.cfg.Model)

	var (
		tts     = cmd.Input.Option("tts")
		ai      = gpt.NewTTSClient(r.cfg.APIKey, r.cfg.Model, tts)
		ctx     = context.Background()
		scanner = bufio.NewScanner(os.Stdin)
		quit    = false
	)

	cmd.PrintText(f2.Apply("\n(use `help` or `quit` to exit.) "))

	for !quit {
		fmt.Print(aurora.Red("\n➜ "))
		if !scanner.Scan() {
			break
		}

		prompt := validateQuestion(scanner.Text())
		switch prompt {
		case "quit":
			quit = true
		case "clear":
			cmd.PrintText("Clearing chat")
			r.session.CurrentChat.Messages = r.loadPersona(cmd, r.session.Persona)
			continue
		case "switch": // ask to save conversation
			keys := make([]string, 0, len(r.session.Personas))
			for k := range r.session.Personas {
				keys = append(keys, k)
			}
			var persona = qh.Ask(
				question.NewChoices("What persona do you want to load?", keys).
					SetMultiselect(false).
					SetMaxAttempts(2),
			)

			r.session.CurrentChat.Messages = r.loadPersona(cmd, persona)
		case "save":
			name := qh.Ask(
				question.NewQuestion("save as: (eg. Alexa)"),
			)
			if _, ok := r.session.Chats[name]; !ok {
				r.session.Chats = make(map[string]*openai.ChatCompletionRequest)
			}
			r.session.Chats[name] = r.session.CurrentChat
			quit = true
		case "print":
			printMessages(r.session.CurrentChat, cmd)
		case "help":
			printHelp(cmd)
		case "":
			continue

		default:
			if strings.HasPrefix(prompt, "system ") {
				cmd.PrintText(f3.Apply("Updating system prompt"))
				r.session.CurrentChat.Messages = append(r.session.CurrentChat.Messages, openai.ChatCompletionMessage{
					Role:    "system",
					Content: strings.Replace(prompt, "system ", "", 1),
				})
				continue
			} else {
				r.session.CurrentChat.Messages = append(r.session.CurrentChat.Messages, openai.ChatCompletionMessage{
					Role:    "user",
					Content: prompt,
				})
			}

			r.printTokenCount(cmd, r.session.CurrentChat.Messages)

			r.session.CurrentChat.Messages, err = ai.ChatCompletion(ctx, r.session.CurrentChat)
			if err != nil {
				cmd.PrintText(f3.Apply(err.Error()))
				continue
			}
			// todo if intent do
			// var answer = r.session.CurrentChat.Messages[len(r.session.CurrentChat.Messages)-1].Content
			//
			// if r.session.Persona == "intent" {
			// 	answer = modules.ReplaceContent(answer)
			// // reset the conversation
			// 	r.session.CurrentChat.Messages = r.loadPersona(cmd, r.session.Persona)
			// }
			// iferr.Warn(ai.TTS(answer))
			continue
		}
	}

	return
}

func (r *Runner) printTokenCount(cmd *console.Script, messages []openai.ChatCompletionMessage) {
	var currentMonth = time.Now().Format("06-Jan")
	today, ok := r.session.Cost[currentMonth]
	if !ok {
		today = new(gpt.CostGaugeMeta)
	}

	gauge, err := gpt.NewCostGauge(r.cfg.Model, entity.BpeEncoderFile, entity.BpeVocabFile)
	if err != nil {
		cmd.PrintError(err.Error())
		return
	}

	meta, err := gauge.CountTokensIn(messages)
	if err != nil {
		cmd.PrintError(err.Error())
		return
	}

	msg := fmt.Sprintf("currently using %d tokens or %.2f percent", meta.NumTokens, meta.Percentage)
	switch true {
	case meta.Percentage > 80:
		cmd.PrintCaution(msg)
	case meta.Percentage > 50:
		cmd.PrintWarning(msg)
	case meta.Percentage > 10:
		cmd.PrintSuccess(msg)
	}

	r.session.Cost[currentMonth] = gauge.UpdateCost(today, meta)
}

func (r *Runner) loadPrevChat(cmd *console.Script, qh *question.Helper) {
	keys := make([]string, 0, len(r.session.Chats))
	for k := range r.session.Chats {
		keys = append(keys, k)
	}
	if len(keys) > 0 {
		chat := qh.Ask(
			question.NewChoices("Previous chats", keys).
				SetMultiselect(false).
				SetMaxAttempts(2),
		)
		previous, ok := r.session.Chats[chat]
		if ok {
			r.session.CurrentChat = previous
		}
	} else {
		cmd.PrintNote("no chats to load")
	}
}

func validateQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "loop", "break", "continue", "exit", "block"}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}
