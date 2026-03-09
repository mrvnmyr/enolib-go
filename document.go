package eno

import (
	"fmt"
	"strconv"
)

type Root interface {
	Key() string
	LineNumber() int
}

type Document struct {
	Elements []Root
}

func (d *Document) Field(key string) (*Field, error) {
	var found *Field
	for _, element := range d.Elements {
		if element.Key() != key {
			continue
		}
		field, ok := element.(*Field)
		if !ok {
			return nil, newError(element.LineNumber(), "A field was expected")
		}
		if found != nil {
			return nil, newError(element.LineNumber(), "Only a single field was expected")
		}
		found = field
	}
	return found, nil
}

func (d *Document) Flag(key string) (*Flag, error) {
	var found *Flag
	for _, element := range d.Elements {
		if element.Key() != key {
			continue
		}
		flag, ok := element.(*Flag)
		if !ok {
			return nil, newError(element.LineNumber(), "A flag was expected")
		}
		if found != nil {
			return nil, newError(element.LineNumber(), "Only a single flag was expected")
		}
		found = flag
	}
	return found, nil
}

func (d *Document) Embed(key string) (*Embed, error) {
	var found *Embed
	for _, element := range d.Elements {
		if element.Key() != key {
			continue
		}
		embed, ok := element.(*Embed)
		if !ok {
			return nil, newError(element.LineNumber(), "An embed was expected")
		}
		if found != nil {
			return nil, newError(element.LineNumber(), "Only a single embed was expected")
		}
		found = embed
	}
	return found, nil
}

type Field struct {
	key        string
	lineNumber int
	content    fieldContentKind
	value      *string
	attributes []*Attribute
	items      []*Item
}

type fieldContentKind int

const (
	fieldContentNone fieldContentKind = iota
	fieldContentValue
	fieldContentAttributes
	fieldContentItems
)

func (f *Field) Key() string     { return f.key }
func (f *Field) LineNumber() int { return f.lineNumber }
func (f *Field) Value() *string  { return f.value }
func (f *Field) Attributes() ([]*Attribute, error) {
	if f.content == fieldContentItems {
		return nil, newError(f.lineNumber, "Attributes expected, Items found")
	}
	if f.content == fieldContentValue {
		return nil, newError(f.lineNumber, "Attributes expected, Value found")
	}
	return f.attributes, nil
}

func (f *Field) Items() ([]*Item, error) {
	if f.content == fieldContentAttributes {
		return nil, newError(f.lineNumber, "Items expected, Attributes found")
	}
	if f.content == fieldContentValue {
		return nil, newError(f.lineNumber, "Items expected, Value found")
	}
	return f.items, nil
}

func (f *Field) RequiredValueString() (string, error) {
	if f.content == fieldContentAttributes {
		return "", newError(f.lineNumber, "Value expected, Attributes found")
	}
	if f.content == fieldContentItems {
		return "", newError(f.lineNumber, "Value expected, Items found")
	}
	if f.value == nil {
		return "", newError(f.lineNumber, "Missing value")
	}
	return *f.value, nil
}

func (f *Field) RequiredValueInt() (int, error) {
	value, err := f.RequiredValueString()
	if err != nil {
		return 0, err
	}
	converted, convErr := strconv.Atoi(value)
	if convErr != nil {
		return 0, wrapError(f.lineNumber, convErr)
	}
	return converted, nil
}

func (f *Field) Attribute(key string) (*Attribute, error) {
	attributes, err := f.Attributes()
	if err != nil {
		return nil, err
	}
	var found *Attribute
	for _, attribute := range attributes {
		if attribute.key != key {
			continue
		}
		if found != nil {
			return nil, newError(f.lineNumber, fmt.Sprintf("Multiple attributes with key %s found", key))
		}
		found = attribute
	}
	if found == nil {
		return nil, newError(f.lineNumber, fmt.Sprintf("Missing attribute %s", key))
	}
	return found, nil
}

type Attribute struct {
	key        string
	lineNumber int
	value      *string
}

func (a *Attribute) Key() string     { return a.key }
func (a *Attribute) LineNumber() int { return a.lineNumber }
func (a *Attribute) Value() *string  { return a.value }
func (a *Attribute) RequiredValueString() (string, error) {
	if a.value == nil {
		return "", newError(a.lineNumber, "Missing value")
	}
	return *a.value, nil
}

type Item struct {
	lineNumber int
	value      *string
}

func (i *Item) LineNumber() int { return i.lineNumber }
func (i *Item) Value() *string  { return i.value }
func (i *Item) RequiredValueString() (string, error) {
	if i.value == nil {
		return "", newError(i.lineNumber, "Value expected")
	}
	return *i.value, nil
}

type Flag struct {
	key        string
	lineNumber int
}

func (f *Flag) Key() string     { return f.key }
func (f *Flag) LineNumber() int { return f.lineNumber }

type Embed struct {
	key            string
	lineNumber     int
	terminatorLine int
	value          *string
}

func (e *Embed) Key() string     { return e.key }
func (e *Embed) LineNumber() int { return e.lineNumber }
func (e *Embed) Value() *string  { return e.value }
func (e *Embed) RequiredValueString() (string, error) {
	if e.value == nil {
		return "", newError(e.lineNumber, "Value expected")
	}
	return *e.value, nil
}
