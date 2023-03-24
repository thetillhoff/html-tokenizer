package validator

import "github.com/thetillhoff/html-tokenizer/pkg/htmlTag"

type NodeBase struct {
	Tag      htmlTag.Tag // self
	children []Node
}
