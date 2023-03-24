package htmlTag

type Tag struct {
	Attributes map[string]string
	Type       TagType
	Name       string
	State      TagState

	StartLine        uint64
	StartChar        uint64
	WhitespacePrefix string // TODO
	WhitespaceSuffix string // TODO
}

func (tag Tag) String() string { // Minimal rendering

	var renderedString string

	if tag.Type == Stringdata {
		renderedString = tag.Attributes["data"]
	} else if tag.Type == Comment {
		renderedString = "<!--" + tag.Attributes["data"] + "-->"
	} else {
		renderedString = "<"
		if tag.State == Close {
			renderedString = renderedString + "/"
		}

		renderedString = renderedString + tag.Name // Add tagName

		for attributeName, attributeValue := range tag.Attributes {
			if attributeValue == "" { // Attribute without a value
				renderedString = renderedString + " " + attributeName
			} else { // Attribute with a value
				renderedString = renderedString + " " + attributeName + "=" + attributeValue
			}
		}

		renderedString = renderedString + ">"
	}

	return renderedString
}
