package model

type EntryType string

const (
	EntryTypeMemo    EntryType = "memo"
	EntryTypePaste   EntryType = "paste"
	EntryTypeWebsite EntryType = "website"
)

type Entry interface {
	Title() string
	Subtitle() string
	Content() string
	Type() EntryType
	Viewed() uint
}
