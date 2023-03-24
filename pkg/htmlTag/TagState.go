package htmlTag

type TagState uint8

const (
	Open  TagState = iota // For <...> tags
	Close                 // For </...> tags
)
