package feedback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Marshal(f Feedback) ([]byte, error) {
	return json.Marshal(&f)
}

func Unmarshal(buf []byte) (f Feedback, err error) {
	err = json.Unmarshal(buf, &f)
	return f, err
}

func PostFeedback(f Feedback, url string) error {
	buf, err := Marshal(f)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to post: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}
