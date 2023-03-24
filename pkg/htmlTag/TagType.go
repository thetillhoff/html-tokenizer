package htmlTag

import (
	"log"
	"strings"
)

type TagType uint8 // 256 possible values

const (
	Normal     TagType = iota // For all normal tags like `html`, `a` etc.
	Stringdata                // For text only, not parsed
	Styledata                 // For style only, not parsed
	Scriptdata                // For script only, not parsed
	Comment                   // For comments only, not parsed
)

func ParseTagType(tag string) TagType {
	log.Println("parsing tag type:", tag)
	tag = strings.ToLower(tag) // convert to lowercase before parsing - it's case insensitive anyway
	switch tag {
	case "style":
		return Styledata
	case "script":
		return Scriptdata
	case "!--":
		return Comment
	default:
		return Normal
	}
}
