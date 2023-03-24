package tokenizer

type Mode uint8

const (
	data           Mode = iota // Generic string data
	tagType                    // Reading the type of a tag
	attributeKey               // Reading the name/key of an attribute
	attributeValue             // Reading the value of an attribute
	styledata                  // Reading the contents of a style tag
	scriptdata                 // Reading the contents of a script tag
	commentdata                // Reading the contents of a comment
)

func (mode Mode) String() string {
	switch mode {
	case data:
		return "data"
	case tagType:
		return "tagType"
	case attributeKey:
		return "attributeKey"
	case attributeValue:
		return "attributeValue"
	case styledata:
		return "styledata"
	case scriptdata:
		return "scriptdata"
	case commentdata:
		return "commentdata"
	default:
		return "unkown mode"
	}
}
