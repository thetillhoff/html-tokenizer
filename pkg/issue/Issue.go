package issue

import "strconv"

type Issue struct {
	Level      IssueLevel
	LineNumber uint64
	CharNumber uint64
	Err        error
}

func (issue Issue) String() string {

	var renderedString string

	switch issue.Level { // Prefix level of issue
	case Critical:
		renderedString = "Critical error"
	case Warning:
		renderedString = "Warning"
	}

	renderedString = renderedString + " on line " + strconv.FormatUint(issue.LineNumber, 10) // Add lineNumber

	renderedString = renderedString + " char " + strconv.FormatUint(issue.CharNumber, 10) // Add charNumber

	renderedString = renderedString + ": " + issue.Err.Error()

	return renderedString
}
