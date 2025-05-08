package blocks

import "encoding/xml"

type Block interface {
	GetType() string
	String() string
}

type XmlRoot struct {
	XMLName xml.Name   `xml:"xml"`
	Blocks  []RawBlock `xml:"block"`
}

type RawBlock struct {
	XMLName xml.Name `xml:"block"`
	Type    string   `xml:"type,attr"`
	Field   string   `xml:"field"`
	Values  []Value  `xml:"value"`
}

type Value struct {
	XMLName xml.Name `xml:"value"`
	Name    string   `xml:"name,attr"`
	Block   RawBlock `xml:"block"`
}

func (r RawBlock) GetType() string {
	return r.Type
}

func (r RawBlock) String() string {
	return "<" + r.Type + ">"
}

type EmptyBlock struct {
	RawBlock
}
