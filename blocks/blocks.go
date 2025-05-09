package blocks

import (
	"encoding/xml"
	"strings"
)

type Block interface {
	GetType() string
	String() string
	Continuous() bool
}

type XmlRoot struct {
	XMLName xml.Name   `xml:"xml"`
	Blocks  []RawBlock `xml:"block"`
}

type RawBlock struct {
	XMLName    xml.Name    `xml:"block"`
	Type       string      `xml:"type,attr"`
	Fields     []Field     `xml:"field"`
	Values     []Value     `xml:"value"`
	Mutation   Mutation    `xml:"mutation"`
	Statements []Statement `xml:"statement"`
	Next       *Next       `xml:"next"`
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

type Mutation struct {
	XMLName    xml.Name    `xml:"mutation"`
	ItemCount  int         `xml:"items,attr"`
	LocalNames []LocalName `xml:"localname"`
	Args       []Arg       `xml:"arg"`
	Key        string      `xml:"key,attr"`
}

type LocalName struct {
	XMLName xml.Name `xml:"localname"`
	Name    string   `xml:"name,attr"`
}

type Statement struct {
	XMLName xml.Name  `xml:"statement"`
	Block   *RawBlock `xml:"block"`
}

type Next struct {
	XMLName xml.Name  `xml:"next"`
	Block   *RawBlock `xml:"block"`
}

type Arg struct {
	Name string `xml:"name,attr"`
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

func (r RawBlock) SingleStatement() Statement {
	return r.Statements[0]
}

func (r RawBlock) String() string {
	return "<" + r.Type + ">"
}

func (r RawBlock) Continuous() bool {
	return true
}

func Pad(block Block) string {
	return " " + strings.Replace(block.String(), "\n", "\n  ", -1) + "\n"
}

func PadBody(blocks []Block) string {
	var builder strings.Builder
	for _, block := range blocks {
		builder.WriteString(Pad(block))
	}
	return builder.String()
}

func JoinBlocks(blocks []Block, delimiter string) string {
	opStrings := make([]string, len(blocks))
	for i, op := range blocks {
		opStrings[i] = op.String()
	}
	return strings.Join(opStrings, delimiter)
}
