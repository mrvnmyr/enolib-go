package eno

import "fmt"

type Locale interface {
	AttributeOutsideField(line int, text string) string
	AttributeWithoutKey(line int, text string) string
	EmbedWithoutKey(line int, text string) string
	EscapeWithoutKey(line int, text string) string
	FieldWithoutKey(line int, text string) string
	InvalidAfterEscape(line int, text string) string
	ItemOutsideField(line int, text string) string
	MixedFieldContent(line int) string
	UnterminatedEmbed(key string, line int) string
	UnterminatedEscapedKey(line int, text string) string
}

type defaultLocale struct{}

func (defaultLocale) AttributeOutsideField(line int, text string) string {
	return fmt.Sprintf("The attribute in line %d is not contained within a field. ('%s')", line, text)
}

func (defaultLocale) AttributeWithoutKey(line int, text string) string {
	return fmt.Sprintf("The attribute in line %d has no key. ('%s')", line, text)
}

func (defaultLocale) EmbedWithoutKey(line int, text string) string {
	return fmt.Sprintf("The embed in line %d has no key. ('%s')", line, text)
}

func (defaultLocale) EscapeWithoutKey(line int, text string) string {
	return fmt.Sprintf("The escape sequence in line %d specifies no key. ('%s')", line, text)
}

func (defaultLocale) FieldWithoutKey(line int, text string) string {
	return fmt.Sprintf("The field in line %d has no key. ('%s')", line, text)
}

func (defaultLocale) InvalidAfterEscape(line int, text string) string {
	return fmt.Sprintf("The escape sequence in line %d can only be followed by an attribute or field operator. ('%s')", line, text)
}

func (defaultLocale) ItemOutsideField(line int, text string) string {
	return fmt.Sprintf("The item in line %d is not contained within a field. ('%s')", line, text)
}

func (defaultLocale) MixedFieldContent(line int) string {
	return fmt.Sprintf("The field in line %d must contain either only attributes, only items, or only a value.", line)
}

func (defaultLocale) UnterminatedEmbed(key string, line int) string {
	return fmt.Sprintf("The embed '%s' starting in line %d is not terminated until the end of the document.", key, line)
}

func (defaultLocale) UnterminatedEscapedKey(line int, text string) string {
	return fmt.Sprintf("The key escape sequence in line %d is not terminated before the end of the line. ('%s')", line, text)
}

var DefaultLocale Locale = defaultLocale{}
