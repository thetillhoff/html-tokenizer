# html-tokenizer

This project is aimed to split html files into a nodetree.
It is meant to be used as a package, and therefore not a standalone project.

## Process

The process looks like this:
- input is a string -> the html document. For now, parsing a stream is not considered
- return value is a list of Nodes and Errors, which is are custom types, and a list of errors for this and all subnodes
- if input string has leading or trailing whitespace (except for trailing newline only), add that into the warning list of the root node


types:
- Node
  - errors []Error
  - lineOffset int // passed into the subnodes on the method calls, so they can calculate their own positions
  - intendation int // passed into subnodes on the method calls, so they can validate themselves and calculate their own positions
  - getErrors(lineOffset int) // Recursive dive into all errors in this node and its children
  - attributes map[string]Attribute
  - children []Node
- NodeA
  - Node
  - parse() // Custom for each Nodetype, 
  - attributeValidations map[string]func(Attribute)
  - ...
- NodeP
  - Node
  - parse()
  - ...
- NodeComment
- ...

- Error
  - text string
  - line int
  - char int
  - level ErrorLevel
- ErrorLevel enum // Possible values are critical and warning (for now)

- Attribute
  - name string
- AttributeInt
  - Attribute
  - Value int
- AttributeString
  - Attribute
  - Value string
- AttributeBool
  - Attribute
  - Value bool


methods:
- validateNode() // Checks the type and its attributes to be valid together, but also whether the types/contents of the attributes fit it

sidenotes:
- the errors are custom types (internally?), so multiple errors can be passed back in one go and also so that line and char numbers can be returned properly

// TODO
- Node.isThisCssSelectorUsed(cssSelector string) bool // Checks if this node or any child nodes fit this css selector
- Node.areAllCssClassesDefined([]definedCssSelector) bool // Checks if any the defined css classes is never used by this node or any child nodes
- validation of
  - indentation
  - allowed child node types
  - closing tags present
  - attribute key names
  - attribute key name combinations
  - attribute value type has to fit attribute key
  - attribute value could be limited to certain values depending on attribute key
- docs for
  - whitespace collapsing on text (leading and trailing whitespace of stringdata is removed, as well as multiple whitespace chars directly after another are all collapsed into a single space, all whitespace is transformed into single spaces)

- todo in tokenizer, all issues are actually errors.
  issue package is unnecessary