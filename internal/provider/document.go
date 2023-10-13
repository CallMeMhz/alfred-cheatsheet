package provider

import "github.com/callmemhz/godash/model"

type Document struct {
	title    string
	subtitle string
	content  string
	typ      model.EntryType
	viewed   uint
}

func (doc *Document) Title() string         { return doc.title }
func (doc *Document) Subtitle() string      { return doc.subtitle }
func (doc *Document) Content() string       { return doc.content }
func (doc *Document) Type() model.EntryType { return doc.typ }
func (doc *Document) Viewed() uint          { return 0 }
