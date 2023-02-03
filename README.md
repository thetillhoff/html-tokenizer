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
- Tokenize(string) ([]Node, []Error) // Main method of this package, returns the tokenized contents of the string
  // Will take the first node in the string and check if it's selfclosing (like <br> or <img>), then either detect its end or call parse() with it and its contents, then does the same with the remaining string
  // Finds the first index of '<', verifies it is on position 0, then takes everything up to the next space, then does a switch case on it to determine the type
  // If the first index of '<' is not 0, verifies that everything before is only whitespace, and if there is a linebreak it in, verifies that all except the last don't have other whitespace in it, else add it to the warnings. If the last linebreak has trailing whitespace, make sure to check every following linebreak to have the same whitespace next to it. If not, add it to the warnings. If there is as much whitespace, remove it on each line, so its considered only once

- validateNode() // Checks the type and its attributes to be valid together, but also whether the types/contents of the attributes fit it

sidenotes:
- the errors are custom types (internally?), so multiple errors can be passed back in one go
- have Node.RecursiveWarnings() and Node.RecursiveErrors() functions that goes into the nodetree and returns all warnings and errors of all nodes
- Node.parse() should also check the indentation of the whole node, f.e. 
- verify indentation 


