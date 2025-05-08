package blocks

import (
	"encoding/xml"
	"strings"
)

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
	Fields  []Field  `xml:"field"`
	Values  []Value  `xml:"value"`
}

type Field struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type Value struct {
	XMLName xml.Name `xml:"value"`
	Name    string   `xml:"name,attr"`
	Block   RawBlock `xml:"block"`
}

type EmptyBlock struct {
	RawBlock
}

func (r RawBlock) GetType() string {
	return r.Type
}

func (r RawBlock) SingleValue() RawBlock {
	return r.Values[0].Block
}

func (r RawBlock) SingleField() string {
	return r.Fields[0].Value
}

func (r RawBlock) String() string {
	return "<" + r.Type + ">"
}

func JoinBlocks(blocks []Block, delimiter string) string {
	opStrings := make([]string, len(blocks))
	for i, op := range blocks {
		opStrings[i] = op.String()
	}
	return strings.Join(opStrings, delimiter)
}
