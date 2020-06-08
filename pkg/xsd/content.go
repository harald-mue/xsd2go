package xsd

import (
	"encoding/xml"
)

type GenericContent interface {
	Attributes() []Attribute
	Elements() []Element
	ContainsText() bool
	compile(*Schema, *Element)
}
type SimpleContent struct {
	XMLName   xml.Name   `xml:"http://www.w3.org/2001/XMLSchema simpleContent"`
	Extension *Extension `xml:"extension"`
}

func (sc *SimpleContent) Attributes() []Attribute {
	if sc.Extension != nil {
		return sc.Extension.Attributes
	}
	return []Attribute{}
}

func (sc *SimpleContent) ContainsText() bool {
	return sc.Extension != nil && sc.Extension.ContainsText()
}

func (sc *SimpleContent) Elements() []Element {
	if sc.Extension != nil {
		return sc.Extension.Elements()
	}
	return []Element{}
}

func (sc *SimpleContent) compile(sch *Schema, parentElement *Element) {
	if sc.Extension != nil {
		sc.Extension.compile(sch, parentElement)
	}
}

type Extension struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2001/XMLSchema extension"`
	Base       reference   `xml:"base,attr"`
	Attributes []Attribute `xml:"attribute"`
	Sequence   *Sequence   `xml:"sequence"`
	typ        Type
}

func (ext *Extension) Elements() []Element {
	elements := []Element{}
	if ext.typ != nil {
		elements = append(elements, ext.typ.Elements()...)
	}
	if ext.Sequence != nil {
		elements = append(elements, ext.Sequence.Elements()...)
		if ext.typ != nil {
			elements = deduplicateElements(elements)
		}
	}
	return elements
}

func deduplicateElements(elements []Element) []Element {
	seen := make(map[string]struct{}, len(elements))
	j := 0
	for _, element := range elements {
		if _, ok := seen[element.GoName()]; ok {
			continue
		}
		seen[element.GoName()] = struct{}{}
		elements[j] = element
		j++
	}
	return elements[:j]
}

func (ext *Extension) ContainsText() bool {
	return ext.Base == "xsd:string" || (ext.typ != nil && ext.typ.ContainsText())
}

func (ext *Extension) compile(sch *Schema, parentElement *Element) {
	if ext.Sequence != nil {
		ext.Sequence.compile(sch, parentElement)
	}
	if ext.Base == "" {
		panic("Not implemented: xsd:extension/@base empty, cannot extend unknown type")
	}

	ext.typ = sch.findReferencedType(ext.Base)
	if ext.typ == nil {
		panic("Cannot build xsd:extension: unknown type: " + string(ext.Base))
	}
	ext.typ.compile(sch, parentElement)
}

type ComplexContent struct {
	XMLName   xml.Name   `xml:"http://www.w3.org/2001/XMLSchema complexContent"`
	Extension *Extension `xml:"extension"`
}

func (cc *ComplexContent) Attributes() []Attribute {
	if cc.Extension != nil {
		return cc.Extension.Attributes
	}
	return []Attribute{}
}

func (cc *ComplexContent) Elements() []Element {
	if cc.Extension != nil {
		return cc.Extension.Elements()
	}
	return []Element{}
}

func (cc *ComplexContent) ContainsText() bool {
	return cc.Extension != nil && cc.Extension.ContainsText()
}

func (c *ComplexContent) compile(sch *Schema, parentElement *Element) {
	if c.Extension != nil {
		c.Extension.compile(sch, parentElement)
	}
}