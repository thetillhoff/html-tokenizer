package tokenizer

import (
	"errors"
	"log"
	"strings"

	tag "github.com/thetillhoff/html-tokenizer/pkg/htmlTag"
	"github.com/thetillhoff/html-tokenizer/pkg/issue"
)

var (
	mode Mode        // Used in Tokenize() and changeMode()
	word string = "" // Used in Tokenize() and changeMode()
)

// Main method of this package, returns the tokenized contents of the string and any issues that occured during the process
func Tokenize(content string) ([]tag.Tag, []issue.Issue) {
	var (
		tags   []tag.Tag
		issues []issue.Issue

		// For attributeKey mode
		constructionTag tag.Tag

		// For attributeValue mode
		attributeName         string
		attributeValueEndRune rune

		lineNumber uint64 = 1
		charNumber uint64 = 1
	)

	mode = data // Start with generic data mode

	for _, char := range content {
		// TODO
		// When parsing the current Node, also add it to the open Nodestack.
		// when a node is a closing Node, make sure to pop open nodes within the Nodestack. (with warning, depending on the nodetype)

		// Global linebreak and location management, so it doesn't have to be done over and over
		if char == '\n' {
			lineNumber = lineNumber + 1 // Increase line number
			charNumber = 1              // Reset char number
		} else {
			charNumber = charNumber + 1 // Incrase char number
		}

		switch mode { // Depending on mode/state, the meaning of chars differentiate
		case data: // Currently no open tag
			switch char {
			case '<':
				if word != "" { // If word is not empty
					// AND not only whitespace // TODO
					constructionTag = tag.Tag{ // Create stringdata tag out of text
						Type:      tag.Stringdata,
						StartLine: lineNumber,
						StartChar: charNumber,
						Attributes: map[string]string{
							"data": strings.TrimSpace(word), // Trim leading and trailing whitespace of word before adding it
						},
					}
					tags = append(tags, constructionTag) // Add tag to list
					log.Println("Adding a tag for stringdata")
				}
				constructionTag = tag.Tag{
					State:      tag.Open,
					StartLine:  lineNumber,
					StartChar:  charNumber,
					Attributes: map[string]string{},
				} // Create new tag that will be constructed by the following chars and set the state to open by default
				mode = changeMode(tagType) // Enter tagType mode
			case '>':
				// Fail, unexpected tag close
				issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: unexpected tag close ('>')")})
				return tags, issues // Syntax error, therefore further parsing not possible
			case ' ':
				if word != "" { // Only if word is not empty, whitespace is important
					if word[len(word)-1] != ' ' { // If last char is already a space
						word = word + string(' ')
					} // Else it will be collapsed / ignored
				} // Else it can be ignored
			case '\t':
				if word != "" { // Only if word is not empty, whitespace is important
					if word[len(word)-1] != ' ' { // If last char is already a space
						word = word + string(' ') // Add space instead of tab
					} // Else it will be collapsed / ignored
				} // Else it can be ignored
			case '\n':
				if word != "" { // Only if word is not empty, whitespace is important
					if word[len(word)-1] != ' ' { // If last char is already a space
						word = word + string(' ') // Add space instead of newline
					} // Else it will be collapsed / ignored
				} // Else it can be ignored
			default:
				word = word + string(char) // Add char to word
			}

		case tagType: // Currently reading the type of a newly opened tag
			switch char {
			case ' ': // tagType now completed
				constructionTag.Type = tag.ParseTagType(word) // tagType now completed, parse it
				constructionTag.Name = word
				switch constructionTag.Type {
				case tag.Comment:
					mode = changeMode(commentdata) // Enter commentdata mode
				default:
					mode = changeMode(attributeKey) // Enter attributeKey mode
				}
			case '\t':
				constructionTag.Type = tag.ParseTagType(word) // tagType now completed, parse it
				constructionTag.Name = word
				switch constructionTag.Type {
				case tag.Comment:
					mode = changeMode(commentdata) // Enter commentdata mode
				default:
					mode = changeMode(attributeKey) // Enter attributeKey mode
				}
			case '\n':
				constructionTag.Type = tag.ParseTagType(word) // tagType now completed, parse it
				constructionTag.Name = word
				switch constructionTag.Type {
				case tag.Comment:
					mode = changeMode(commentdata) // Enter commentdata mode
				default:
					mode = changeMode(attributeKey) // Enter attributeKey mode
				}
			case '>':
				constructionTag.Type = tag.ParseTagType(word) // tagType now completed, parse it
				constructionTag.Name = word
				switch word {
				case "style":
					mode = changeMode(styledata) // Enter style mode
					log.Println("Adding a tag of type", constructionTag.Name)
					tags = append(tags, constructionTag) // Add finished tag to list
					constructionTag = tag.Tag{
						Type:       tag.Styledata,
						StartLine:  lineNumber,
						StartChar:  charNumber,
						Attributes: map[string]string{"data": ""},
					} // Create styledata tag
				case "script":
					mode = changeMode(scriptdata) // Enter script mode
					log.Println("Adding a tag of type", constructionTag.Name)
					tags = append(tags, constructionTag) // Add finished tag to list
					constructionTag = tag.Tag{
						Type:       tag.Styledata,
						StartLine:  lineNumber,
						StartChar:  charNumber,
						Attributes: map[string]string{"data": ""},
					} // Create scriptdata tag
				default:
					mode = changeMode(data) // Enter data mode
					log.Println("Adding a tag of type", constructionTag.Name)
					tags = append(tags, constructionTag) // Add finished tag to list
				}
			case '/':
				if constructionTag.State == tag.Close { // If tag was already identified as a closing tag
					issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: tag was already closed (duplicate '/')")}) // Fail, because two '/' are not allowed in a tag name
					return tags, issues                                                                                                                                                                  // Syntax error, therefore further parsing not possible
				}
				constructionTag.State = tag.Close // Tag is identified as a closing tag
			default:
				word = word + string(char) // Add char to word
			}

		case attributeKey: // Currently reading the name of an attribute of the currently open tag
			switch char {
			case ' ': // TODO same needs to be done for any whitespace in this mode
				if word != "" { // If word is not empty
					constructionTag.Attributes[word] = "" // Add current word to map of attributes of currently open tag
					attributeName = word                  // Set global attributeName variable to current word so it can be found in the following (for cases like `attr = value`)
				}
			case '>':
				if word != "" { // If word is not empty
					constructionTag.Attributes[word] = "" // Add current word to map of attributes of currently open tag
					attributeName = ""                    // Reset attributeName
				}
				switch constructionTag.Name {
				case "style":
					mode = changeMode(styledata) // Enter style mode
					log.Println("Adding a tag of type", constructionTag.Name)
					tags = append(tags, constructionTag) // Add finished tag to list
					constructionTag = tag.Tag{
						Type:       tag.Styledata,
						StartLine:  lineNumber,
						StartChar:  charNumber,
						Attributes: map[string]string{"data": ""},
					} // Create styledata tag
				case "script":
					mode = changeMode(scriptdata) // Enter script mode
					log.Println("Adding a tag of type", constructionTag.Name)
					tags = append(tags, constructionTag) // Add finished tag to list
					constructionTag = tag.Tag{
						Type:       tag.Styledata,
						StartLine:  lineNumber,
						StartChar:  charNumber,
						Attributes: map[string]string{"data": ""},
					} // Create scriptdata tag
				default:
					mode = changeMode(data) // Enter data mode
					log.Println("Adding a tag of type", constructionTag.Name)
					tags = append(tags, constructionTag) // Add finished tag to list
				}
			case '=':
				attributeName = word              // Set global attributeName variable to current word so it can be found in the following (for cases like `attr = value`)
				mode = changeMode(attributeValue) // Enter attributeValueName mode
				attributeValueEndRune = ' '       // Set attributeValueEndingRune to default of space
			case '/':
				issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: attributeKey must not contain '/'")}) // Fail, because '/' is contained in attributeKey
				return tags, issues                                                                                                                                                             // Syntax error, therefore further parsing not possible
			default:
				word = word + string(char) // Add char to word
			}

		case attributeValue: // Currently reading the value of an attribute of the currently open tag
			switch char {
			case '\'':
				if word == "" { // No value yet
					attributeValueEndRune = '\'' // Set attributeValueEndRune to single quote
					word = word + string(char)   // Add char to word
				} else if char == attributeValueEndRune { // End value
					word = word + string(char)                       // Add char to word
					constructionTag.Attributes[attributeName] = word // Save current word as attributeValue to current attribute
					attributeName = ""                               // Reset attributeName
					mode = changeMode(attributeKey)                  // Enter attributeKey mode
				}
			case '"':
				if word == "" { // No value yet
					attributeValueEndRune = '"' // Set attributeValueEndRune to double quotes
					word = word + string(char)  // Add char to word
				} else if char == attributeValueEndRune { // End value
					word = word + string(char)                       // Add char to word
					constructionTag.Attributes[attributeName] = word // Save current word as attributeValue to current attribute
					attributeName = ""                               // Reset attributeName
					mode = changeMode(attributeKey)                  // Enter attributeKey mode
				}
			case ' ':
				if char == attributeValueEndRune { // End value & attribute
					if word != "" { // If word has already started
						constructionTag.Attributes[attributeName] = word // Save current word as attributeValue to current attribute
						attributeName = ""                               // Reset attributeName
						mode = changeMode(attributeKey)                  // Enter attributeKey mode
					} // Else, word hasn't started yet, so ignore whitespace before it
				} else { // Just a normal char
					word = word + string(char) // Add char to word
				}
			default:
				word = word + string(char) // Add char to word
			}

		case styledata:
			word = word + string(char) // Add char to word
			if strings.HasSuffix(word, "</style>") {
				constructionTag.Attributes["data"] = strings.TrimSuffix(word, "</style>") // Add contents of style node as data
				constructionTag = tag.Tag{                                                // Create styledata tag out of text
					Type:      tag.Styledata,
					StartLine: lineNumber,
					StartChar: charNumber,
					Attributes: map[string]string{
						"data": word,
					},
				}

				mode = changeMode(data)
				log.Println("Adding a tag of type", constructionTag.Name)
				tags = append(tags, constructionTag) // Add tag to list
			}
		case scriptdata:
			word = word + string(char) // Add char to word
			if strings.HasSuffix(word, "</script>") {
				constructionTag.Attributes["data"] = strings.TrimSuffix(word, "</script>") // Add contents of script node as data

				mode = changeMode(data)
				log.Println("Adding a tag of type", constructionTag.Name)
				tags = append(tags, constructionTag) // Add tag to list
			}
		case commentdata:
			word = word + string(char) // Add char to word
			if strings.HasSuffix(word, "-->") {
				constructionTag.Attributes["data"] = strings.TrimSuffix(word, "-->") // Add contents of comment as data

				mode = changeMode(data)
				log.Println("Adding a tag of type", constructionTag.Name)
				tags = append(tags, constructionTag) // Add tag to list
			}
		}
	}

	switch mode { // Validate mode after whole stream was read (has to be in data mode)
	case data: // Only valid state
		if word != "" { // If there is stringdata trailing without a tag around it, it's a synthax error
			issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: trailing stringdata is not allowed (end mode: stringdata)")}) // Fail, because
			return tags, issues                                                                                                                                                                                     // Syntax error, therefore further parsing not possible
		}
	case tagType: // Invalid state
		issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: last tag wasn't closed (end mode: tagType)")}) // Fail, because
		return tags, issues                                                                                                                                                                      // Syntax error, therefore further parsing not possible
	case attributeKey: // Invalid state
		issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: last tag wasn't closed (end mode: attributeKey)")}) // Fail, because
		return tags, issues                                                                                                                                                                           // Syntax error, therefore further parsing not possible
	case attributeValue: // Invalid state
		issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: last tag wasn't closed (end mode: attributeValue)")}) // Fail, because
		return tags, issues                                                                                                                                                                             // Syntax error, therefore further parsing not possible
	case styledata: // Invalid state
		issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: style tag wasn't closed (end mode: styledata)")}) // Fail, because
		return tags, issues                                                                                                                                                                         // Syntax error, therefore further parsing not possible
	case scriptdata: // Invalid state
		issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: script tag wasn't closed (end mode: scriptdata)")}) // Fail, because
		return tags, issues                                                                                                                                                                           // Syntax error, therefore further parsing not possible
	case commentdata: // Invalid state
		issues = append(issues, issue.Issue{Level: issue.Critical, LineNumber: lineNumber, CharNumber: charNumber, Err: errors.New("syntax error: comment tag wasn't closed (end mode: commentdata)")}) // Fail, because
		return tags, issues                                                                                                                                                                             // Syntax error, therefore further parsing not possible
	}

	return tags, issues
}
