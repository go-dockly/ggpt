package runner

import (
	console "github.com/DrSmithFr/go-console"
	"github.com/sashabaranov/go-openai"
)

func printMessages(chat *openai.ChatCompletionRequest, cmd *console.Script) {
	for _, v := range chat.Messages {
		if v.Role == "assistant" {
			cmd.PrintText(f3.Apply(v.Role + ": " + v.Content))
		} else if v.Role == "user" {
			cmd.PrintText(f1.Apply(v.Role + ": " + v.Content))
		} else {
			cmd.PrintText(v.Role + ": " + v.Content)
		}
	}
}

func printHelp(cmd *console.Script) {
	cmd.PrintText(`
    Keyword usage:
		clear (current conversation) 
		save (current conversation) 
		switch (start new chat. optionally load a different persona)
		system <prompt> (add system instruction prompt to current conversation)
		quit (exit the application) 
`)
}
