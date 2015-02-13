package feedback

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var Faces = []string{
	"x0",
	">8c",
	":c",
	":(",
	":|",
	":)",
	":D",
	"<8D",
}

var FaceMsg = `
	0	x0	horrible
	1	>8c	angering
	2	:c	very sad
	3	:(	sad
	4	:|	meh
	5	:)	happy
	6	:D	very happy
	7	<8D	fantastic
`

type Feedback struct {
	Name      string
	Email     string
	Message   string
	Questions []string
	Answers   []string
	ScoreFace string
	ScoreInt  int
}

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

var (
	ScorePrompt     = "Score your experience, "
	ScoreIntPrompt  = "from 0 (horrible) to 7 (fantastic!!!)"
	ScoreFacePrompt = "with a number:" + FaceMsg
)

func PromptForFeedback(rw io.ReadWriter, opt Options) (Feedback, error) {
	var f Feedback
	var err error

	// input validation
	if len(opt.Questions) != len(opt.QuestionsRequired) {
		return f, errors.New("Questions and QuestionsRequired length mismatch")
	}
	for i, q := range opt.Questions {
		if q == "" {
			return f, fmt.Errorf("empty question %d", i)
		}
	}

	// internal prompt func for niceness below
	prompt := func(prompt, defaultPrompt string, skip, required bool) (ans string) {
		if err != nil {
			return "" // we're bailing.
		}

		if skip {
			return ""
		}

		if prompt == "" {
			prompt = defaultPrompt
		}

		ans, err = promptOnce(rw, prompt, required)
		return ans
	}

	f.Name = prompt(opt.NamePrompt, "Name", opt.NameSkip, opt.NameRequired)
	f.Email = prompt(opt.EmailPrompt, "Email", opt.EmailSkip, opt.EmailRequired)
	f.Message = prompt(opt.MessagePrompt, "Message", opt.MessageSkip, opt.MessageRequired)

	// questions
	f.Questions = opt.Questions
	f.Answers = make([]string, len(f.Questions))
	for i, q := range opt.Questions {
		f.Answers[i] = prompt(q, "<question>", false, opt.QuestionsRequired[i])
	}

	// score
	var scorePrompt = ScorePrompt
	if opt.ScoreFaces {
		scorePrompt += ScoreFacePrompt
	} else {
		scorePrompt += ScoreIntPrompt
	}
	for err == nil {
		ans := prompt(scorePrompt, "<score>", opt.ScoreSkip, opt.ScoreRequired)
		if i, ierr := strconv.Atoi(ans); ierr == nil {
			f.ScoreInt = i
			f.ScoreFace = Faces[i]
			break
		}
		if ans == "" && !opt.ScoreRequired {
			f.ScoreInt = -1
			f.ScoreFace = "N/A"
			break
		}
		if ans != "" {
			fmt.Fprintf(rw, "invalid answer: %s\n", ans)
		}
	}
	return f, err
}

func promptOnce(rw io.ReadWriter, prompt string, required bool) (answer string, err error) {
	r := ""
	if required {
		r = " (required)"
	}

	br := bufio.NewReader(rw)

	for {
		if _, err = fmt.Fprintf(rw, "%s%s: ", prompt, r); err != nil {
			return "", err
		}

		if answer, err = br.ReadString('\n'); err != nil {
			return "", err
		}
		answer = answer[:len(answer)-1] // remove '\n'
		if answer == "" && required {
			continue
		}
		return answer, nil
	}
}
