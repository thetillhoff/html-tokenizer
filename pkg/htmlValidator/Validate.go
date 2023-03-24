package validator

import (
	"github.com/thetillhoff/html-tokenizer/pkg/htmlTag"
	"github.com/thetillhoff/html-tokenizer/pkg/issue"
)

func Validate(tags []htmlTag.Tag) ([]Node, []issue.Issue) {
	var (
		rootNodes = []Node{}
		issues    = []issue.Issue{}
	)

	// TODO implement
	// - different node types
	// - allowed childen types
	// - allowed attribute keys
	// - allowed attribute values per key

	return rootNodes, issues
}
