package format

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/pkg/errors"
)

const (
	startTimerMsg = `{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Time to focus! You've got *{{ .TimeDuration }}* minutes."
			}
		}
	]
}`

	pauseTimerMsg = `{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Timer is paused. Type /gomato_resume to resume the timer or /gomato_stop to delete the timer."
			}
		}
	]
}`

	resumeTimerMsg = `{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Resuming the timer!"
			}
		}
	]
}`

	stopTimerMsg = `{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "The time has stopped. Use /gomato_start to start a new timer."
			}
		}
	]
}`
)

// Message formats a message based on the given scenario
func Message(messageType string, data interface{}) (interface{}, error) {
	var text string
	switch messageType {
	case "start":
		text = startTimerMsg
	case "pause":
		text = pauseTimerMsg
	case "resume":
		text = resumeTimerMsg
	case "stop":
		text = stopTimerMsg
	default:
		return nil, errors.New("invalid message type")
	}
	var buf bytes.Buffer
	tmpl, err := template.New("message").Parse(text)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse message template")
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, errors.Wrap(err, "failed to execute message template")
	}

	var resp interface{}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal template json")
	}

	return resp, nil
}
