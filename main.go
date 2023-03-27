package main

import (
	console "github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/input/argument"
	"github.com/DrSmithFr/go-console/input/option"
	"github.com/go-dockly/ggpt/pkg/runner"
	"github.com/go-dockly/utility/xerrors/iferr"
)

func main() {

	r, err := runner.New()
	iferr.Exit(err)

	script := console.Command{
		UseNamespace: true,
		Description:  "ggpt: command line tool",
		Scripts: []*console.Script{
			{
				Name:        "q",
				Description: "Ask a question to gpt.",
				Arguments: []console.Argument{
					{
						Name:        "query",
						Description: "eg: Generate git commit message, changes: $(git diff)",
						Value:       argument.Required,
					},
				},
				Options: []console.Option{
					{
						Name:        "persona",
						Shortcut:    "p",
						Value:       option.Optional,
						Description: "Load custom persona for gpt.",
					},
				},
				Runner: r.Completion,
			},
			{
				Name:        "voice",
				Description: "Use microphone to interact with gpt.",
				Arguments: []console.Argument{
					{
						Name:        "persona",
						Description: "Load custom persona for gpt.",
						Value:       argument.Optional,
					},
				},
				Options: []console.Option{
					{
						Name:        "tts",
						Shortcut:    "t",
						Value:       option.Optional,
						Description: "Use text to speech.",
					},
				},
				Runner: r.Recorder,
			},
			{
				Name:        "play",
				Description: "Play a wav audio file.",
				Arguments: []console.Argument{
					{
						Name:        "file",
						Description: "Path to file eg: sound.wav",
						Value:       argument.Required,
					},
				},
				Runner: r.PlaySound,
			},
			{
				Name:        "transcribe",
				Description: "Transcribe an audio file with whisper (wav or mp3)",
				Arguments: []console.Argument{
					{
						Name:        "file",
						Description: "Path to file eg: speech.wav",
						Value:       argument.Required,
					},
				},
				Runner: r.TranscribeFile,
			},
			{
				Name:        "say",
				Description: "Text to speech. Requires { docker run -it -p 5002:5002 synesthesiam/mozillatts }",
				Arguments: []console.Argument{
					{
						Name:        "message",
						Description: "Say this message.",
						Value:       argument.Required,
					},
				},
				Runner: r.Say,
			},
			{
				Name:        "ichat",
				Description: "Start an interactive chat session with gpt.",
				Arguments: []console.Argument{
					{
						Name:        "persona",
						Description: "Load custom persona for gpt.",
						Value:       argument.Optional,
					},
				},
				Options: []console.Option{
					{
						Name:        "tts",
						Shortcut:    "t",
						Value:       option.Optional,
						Description: "Use text to speech.",
					},
				},
				Runner: r.IChat,
			},
			{
				Name:        "chat",
				Description: "Start a non-interactive chat with gpt.",
				Arguments: []console.Argument{
					{
						Name:        "name",
						Description: "eg: number",
						Value:       argument.Required,
					},
				},
				Options: []console.Option{
					{
						Name:        "tts",
						Shortcut:    "t",
						Value:       option.Optional,
						Description: "Use text to speech.",
					},
					{
						Name:        "user",
						Shortcut:    "u",
						Value:       option.Optional,
						Description: "User prompt to gpt.",
					},
					{
						Name:        "system",
						Shortcut:    "s",
						Value:       option.Optional,
						Description: "Gpt system prompt injection.",
					},
				},
				Runner: r.Chat,
			},
			{
				Name:        "delete-chat",
				Description: "Remove a stored chat.",
				Arguments: []console.Argument{
					{
						Name:        "name",
						Description: "eg: number",
						Value:       argument.Required,
					},
				},
				Runner: r.DeleteChat,
			},
			{
				Name:        "chats",
				Description: "List all conversations or select a previously started chat by name.",
				Arguments: []console.Argument{
					{
						Name:        "name",
						Description: "Load a previous conversation by name.",
						Value:       argument.Optional,
					},
				},
				Options: []console.Option{
					{
						Name:        "tts",
						Shortcut:    "t",
						Value:       option.None,
						Description: "Use text to speech.",
					},
				},
				Runner: r.Chats,
			},
			{
				Name:        "show-chat",
				Description: "Print messages from a saved conversation.",
				Arguments: []console.Argument{
					{
						Name:        "name",
						Description: "Name of chat to load.",
						Value:       argument.Optional,
					},
				},
				Runner: r.ShowChat,
			},
			{
				Name:        "show-persona",
				Description: "Print profile of custom gpt persona.",
				Arguments: []console.Argument{
					{
						Name:        "persona",
						Description: "Print persona profile content prompt.",
						Value:       argument.Optional,
					},
				},
				Runner: r.PrintPersonas,
			},
			{
				Name:        "count",
				Description: "Count tokens in a file.",
				Arguments: []console.Argument{
					{
						Name:        "file",
						Description: "Path to file eg: main.go",
						Value:       argument.Required,
					},
				},
				Runner: r.CountTokens,
			},
			{
				Name:        "model",
				Description: "Print or update current gpt model used in future conversations.",
				Arguments: []console.Argument{
					{
						Name:        "model",
						Description: "eg: gpt-4 omit to print current model.",
						Value:       argument.Optional,
					},
				},
				Runner: r.SetModel,
			},
			{
				Name:        "encode",
				Description: "Encode a word using the bpe tokenizer.",
				Arguments: []console.Argument{
					{
						Name:        "word",
						Description: "eg: Paris -> 6342",
						Value:       argument.Optional,
					},
				},
				Runner: r.Encode,
			},
			{
				Name:        "show-cost",
				Description: "Show cost of ggpt this month.",
				Arguments: []console.Argument{
					{
						Name:        "month",
						Description: "show cost for specific month eg: 23-Mar",
						Value:       argument.Optional,
					},
				},
				Runner: r.Cost,
			},
			{
				Name:        "sql",
				Description: "Translate natural language to SQL query.",
				Arguments: []console.Argument{
					{
						Name:        "text",
						Description: "eg: find all the cats in the box",
						Value:       argument.Required,
					},
					{
						Name:        "schema",
						Description: "eg: data/schema.sql",
						Value:       argument.Optional,
					},
				},
				Runner: r.TranslateToSQLQuery,
			},
			{
				Name:        "sqlnl",
				Description: "Translate SQL query to natural language.",
				Arguments: []console.Argument{
					{
						Name:        "query",
						Description: "eg: SELECT * FROM cats WHERE color = 'grey';",
						Value:       argument.Required,
					},
				},
				Runner: r.TranslateSQLQuery,
			},
			{
				Name:        "shell",
				Description: "Explain shell command.",
				Arguments: []console.Argument{
					{
						Name:        "command",
						Description: "Translate shell command.",
						Value:       argument.Required,
					},
				},
				Runner: r.ExplainShellCommand,
			},
			{
				Name:        "se",
				Description: "Execute shell command.",
				Arguments: []console.Argument{
					{
						Name:        "command",
						Description: "Run shell command.",
						Value:       argument.Required,
					},
				},
				Runner: r.ExecuteShellCommand,
			},
			{
				Name:        "regex",
				Description: "Translate regular expression.",
				Arguments: []console.Argument{
					{
						Name:        "expression",
						Description: "Translate regular expression.",
						Value:       argument.Required,
					},
				},
				Runner: r.TranslateRegex,
			},
			{
				Name:        "re",
				Description: "Explain regular expression.",
				Arguments: []console.Argument{
					{
						Name:        "expression",
						Description: "Explain regular expression.",
						Value:       argument.Required,
					},
				},
				Runner: r.ExplainRegex,
			},
			{
				Name:        "settings",
				Description: "Update gpt settings.",
				Options: []console.Option{
					{
						Name:        "api_key",
						Shortcut:    "a",
						Value:       option.Optional,
						Description: "Set api key for gpt.",
					},
					{
						Name:        "user",
						Shortcut:    "u",
						Value:       option.Optional,
						Description: "Update user name.",
					},
					{
						Name:        "profile",
						Shortcut:    "p",
						Value:       option.Optional,
						Description: "eg apply `creative` or `focused` settings",
					},
					{
						Name:        "temperature",
						Shortcut:    "t",
						Value:       option.Optional,
						Description: "Set temperature to use with gpt model.",
					},
					{
						Name:        "max_tokens",
						Shortcut:    "mt",
						Value:       option.Optional,
						Description: "Limit the tokens spend on a completion.",
					},
					{
						Name:        "top_p",
						Shortcut:    "tp",
						Value:       option.Optional,
						Description: "range from 0 to 1 (0.5 will make the generated text more focused, while closer to 1 more diverse)",
					},
					{
						Name:        "n_completions",
						Shortcut:    "nc",
						Value:       option.Optional,
						Description: "How many example prompts should gpt generate.",
					},
					{
						Name:        "presence_penalty",
						Shortcut:    "pp",
						Value:       option.Optional,
						Description: "range from -2 to 2, positive values will generate more diverse, less repetitive text.",
					},
					{
						Name:        "frequency_penalty",
						Shortcut:    "fp",
						Value:       option.Optional,
						Description: "range from -2 to 2, negative values will make the model use more frequent tokens.",
					},
					{
						Name:        "logit_bias",
						Shortcut:    "l",
						Value:       option.Optional,
						Description: "Avoid or Emphasize certain words with positive or negative bias (eg avoid Paris:-10)",
					},
					{
						Name:        "stop",
						Shortcut:    "s",
						Value:       option.Optional,
						Description: "Set a list of stop sequences.",
					},
				},
				Runner: r.UpdateSettings,
			},
		},
	}

	script.Run()
}
