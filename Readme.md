
# GoGPT
A command-line productivity tool powered by OpenAI's ChatGPT. Leverage ChatGPT's capabilities to generate shell commands, code snippets, comments, documentation and more. Forget about cheat sheets and notes, with this tool you can get accurate answers right in your terminal. Reduce your daily Google searches and save valuable time.

Programmable command-line voice chat interface with text to speech output

## Prerequisites

Some functionalities like audio requires additional libraries to be installed
```shell
➜  brew install swig  -> snowboy hotword detection dependency
➜  brew install portaudio -> recording microphone
```

If you want to take mozilla/tacotron2 for a spin for TTS do the following in preparation:
```shell
➜  docker run -it -p 5002:5002 synesthesiam/mozillatts
```

Test on mac by replacing (aplay -) with (ffplay -autoexit -nodisp -)
```shell
➜  curl -G --output - \
    --data-urlencode 'text=Welcome to the world of speech!' \
    'http://localhost:5002/api/tts' | \
    aplay -
```

## Installation
```shell
➜  go install github.com/go-dockly/ggpt
```
You'll need an OpenAI API key, you can generate one [here](https://beta.openai.com/account/api-keys).

You will be prompted for your key to be stored in `~/.ggpt/config.json`.

## Usage
`ggpt` has a variety of use cases, including simple queries, shell comands, code snippets, etc.
### Simple queries
Use it as a search engine:
```shell
➜  ggpt q "nginx default config file location"
# -> The default configuration file for Nginx is located at /etc/nginx/nginx.conf.
```
```shell
➜  ggpt q "docker show all local images"
# -> You can view all locally available Docker images by running: `docker images`
```
```shell
➜  ggpt q "mass of sun"
# -> = 1.99 × 10^30 kg
```
### Conversion
Convert various units and measurements such as time, distance, weight, temperature, etc.
```shell
➜  ggpt q "1 hour and 30 minutes to seconds"
# -> 5,400 seconds
```
```shell
➜  ggpt q "1 kilometer to mile"
# -> 1 kilometer is equal to 0.62137 miles.
```
```shell
➜  ggpt q "$(date) to Unix timestamp"
# -> The Unix timestamp for Thu Mar 2 00:13:11 CET 2023 is 1677327191.
```
### SQL natural language interface
Translate natural language to SQL query.
```shell
➜  ggpt sql "show me all the cars that are red" "data/schema.sql"
# -> SELECT * FROM cars WHERE color = "red";
```
Alternatively, you can use a custom schema to improve the prompt
```shell
➜  ggpt sql "show me all the cars that are red" "data/schema.sql"
# -> Sorry, the table schema provided is for "cats" and not for "cars". Please provide the correct table schema to proceed. 
```
Translate an SQL query into natural language.
```shell
➜  ggpt sqlnl "SELECT * FROM cats WHERE color = 'grey'"
# -> Retrieve all information from the table "cats" where the color is "grey".
```
### Shell commands
Look up the syntax of any shell command with the `shell` argument, quickly find and execute the commands you need right in the terminal with the short-cut `se`.

```shell
➜  ggpt shell "change all files in current directory to read only"
# -> chmod 444 *
➜  ggpt se "make all files in current directory read only"
# -> chmod 444 *
# -> Execute shell command? [y/N]: y
# ...
```
Go GPT is aware of OS and `$SHELL`, it will produce shell command according to your system. Try ask `ggpt` to update your system. Here's an example using macOS:
```shell
➜  ggpt se "update my system"
# -> sudo softwareupdate -i -a
```
The same prompt, when used on Ubuntu, will produce:
```shell
➜  ggpt se "update my system"
# -> sudo apt update && sudo apt upgrade -y
```

Let's try starting a docker container:
```shell
➜  ggpt se "start nginx using docker, forward 443 and 80 port, mount current folder with index.html"
# -> docker run -d -p 443:443 -p 80:80 -v $(pwd):/usr/share/nginx/html nginx
# -> Execute shell command? [y/N]: y
# ...
```
Next try passing arguments like output file name to ffmpeg:
```shell
➜  ggpt se "slow down video twice using ffmpeg, input video name \"input.mp4\" output video name \"output.mp4\""
# -> ffmpeg -i input.mp4 -filter:v "setpts=2.0*PTS" output.mp4
# -> Execute shell command? [y/N]: y
# ...
```
GPT magic passing file names to ffmpeg:
```shell
➜  ls
# -> 1.mp4 2.mp4 3.mp4
➜  ggpt se "using ffmpeg combine multiple videos into one without audio. Video file names: $(ls -m)"
# -> ffmpeg -i 1.mp4 -i 2.mp4 -i 3.mp4 -filter_complex "[0:v] [1:v] [2:v] concat=n=3:v=1 [v]" -map "[v]" out.mp4
# -> Execute shell command? [y/N]: y
# ...
```
Ask chatGPT to generate your git commit message:
```shell
➜  ggpt q "Generate git commit message, my changes: $(git diff)"
# -> Commit message: Implement Model enum and getPrompt() func, add temperature ands top_p args for OpenAI request.
```
Ask chatGPT to investigate logs from your docker container directly
```shell
➜  ggpt q "check these logs, find errors, and explain what the error is about: ${docker logs -n 20 container_name}"
# ...
```
### Regular expressions
Look up the syntax of any regex with the `regex` argument, quickly find the explanation with the short-cut `re`.

```shell
➜  ggpt regex "match a string without the string hello"
# -> ^(?!.*hello).*$
➜  ggpt re "^(?!.*hello).*$"
# -> The expression uses a negative lookahead (?!hello) to ensure that the string does not contain the word "hello".
# -> The "^" character matches the beginning of the string, and the ".*" matches any character zero or more times. 
# -> The "$" character matches the end of the string.
# ...
```
### Generating code
With `code` parameters we can query only code as output, for example:
```shell
➜  ggpt code "Solve classic fizz buzz problem using Golang"
```
```golang
package main

import "fmt"

func main() {
    for i := 1; i <= 100; i++ {
        if i%3 == 0 && i%5 == 0 {
            fmt.Println("FizzBuzz")
        } else if i%3 == 0 {
            fmt.Println("Fizz")
        } else if i%5 == 0 {
            fmt.Println("Buzz")
        } else {
            fmt.Println(i)
        }
    }
}
```
Since it is valid golang code, we can redirect the output to a file:
```shell
➜  ggpt code "solve classic fizz buzz problem using Golang without comments" > fizz_buzz.go
➜  go run fizz_buzz.go
# 1
# 2
# Fizz
# 4
# Buzz
# Fizz
# ...
```

### Chat 
if chat does not exist new chat will be created
<!-- ggpt chat number -u/--user "remember 3" -s/--system "you are a calculator"-->
To start or contine a chat session, use the `chat` argument followed by a unique session name and a prompt:
```shell
➜  ggpt chat number -p "please remember my favorite number: 4"
# -> I will remember that your favorite number is 4.
➜  ggpt chat number -p "to my favorite number add 4?"
# -> Your favorite number is 4, so if we add 4 to it, the result would be 8.
➜  ggpt rm-chat number
# -> delete chat number
```

### Chat sessions

Feel like having an interactive chat session with gpt, use the `chat` argument:
```shell
➜  ggpt ichat glados --tts=true
# Loading profile for glados
# Using gpt model gpt-3.5-turbo

# (type `quit` to exit or chat away) 

# ➜ what is the size of jupiter?
# Who cares about the size of Jupiter.
```

If you have a microphone connected, use the `voice` argument instead with snowboy hotword detection for 'computer':
```shell
➜  ggpt voice glados --tts=true
# Loading profile for glados
# Using gpt model gpt-3.5-turbo
# microphone streaming: sample_rate=16000, bit_depth=16, hotword=computer, silence_counter=1.5sec

# ➜ what is the size of jupiter?
# Who cares about the size of Jupiter.
```

To list all the current chat sessions, use the `chats` option:
```shell
➜  ggpt chats
# .../ggpt/chats/number
```
To show all the messages related to a specific chat session, use the `show-chat` option followed by the session name:
```shell
➜  ggpt show-chat number
# user: please remember my favorite number: 4
# assistant: I will remember that your favorite number is 4.
# user: what would be my favorite number + 4?
# assistant: Your favorite number is 4, so if we add 4 to it, the result would be 8.
```

### Settings control
Wizard to change settings used against openAI such as temperature, etc and save to config
```shell
➜  ggpt settings -t/--temperature 0.7 -l/--logit_bias "Paris:-10"
# https://algowriting.medium.com/gpt-3-temperature-setting-101-41200ff0d0be
```
Alternatively load a preset settings profile
```shell
➜  ggpt settings -p/--profile creative/focused/default
```

### Cost analysis
Show cost of ggpt for a specific month (default: current).
```shell
➜  ggpt cost "23-Mar"
# cost tokens: $0.01
# total tokens processed 1698
# cost whisper $0.02
# total seconds transcribed 0.60
```

