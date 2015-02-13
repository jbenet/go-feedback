# go get feedback

`feedback` is a small pkg and cli toolset to collect user feedback.
It's super trivial.

## Install

```go
go get github.com/jbenet/go-feedback
```

## Usage

### lib

```go
import (
  feedback "github.com/jbenet/go-feedback"
)

func main() {
  // this really should be in the io pkg...
  rw := feedback.NewReadWriter(os.Stdin, os.Stderr)

  // prompt the user for feedback through a ReadWriter.
  // this could be stdin/stdout or a net.Conn or whatever
  f, err := feedback.PromptForFeedback(rw, feedback.Options{})

  // marshal it into json & write it out
  buf, err := feedback.Marshal(f)
  os.Stdout.Write(buf)

  // or POST it via http
  feedback.Post(f, "http://feedback.io")
}
```

Options:

```go
type Options struct {
  NameSkip     bool
  NamePrompt   string
  NameRequired bool

  EmailSkip     bool
  EmailPrompt   string
  EmailRequired bool

  MessageSkip     bool
  MessagePrompt   string
  MessageRequired bool

  Questions         []string
  QuestionsRequired []bool

  ScoreSkip     bool
  ScorePrompt   string
  ScoreRequired bool
  ScoreFaces    bool
}
```

### cli

Client

```sh
> go get github.com/jbenet/go-feedback/feedback
> feedback >result                        # output to stdout
> feedback --post http://127.0.0.1:1234   # or post to url
Name: jbenet
Email: juan@benet.ai
Message: This is sweet
Score your experience, from 0 (horrible) to 7 (fantastic!!!): 7
```

Server

```sh
> go get github.com/jbenet/go-feedback/feedback-server
> feedback-server -l :1234
Listening at 0.0.0.0:1234
{"Name":"","Email":"","Message":"This is sweet","Questions":null,"Answers":[],"ScoreFace":"\u003c8D","ScoreInt":7}
```
