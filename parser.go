package eno

import "strings"

func Parse(input string) (*Document, error) {
	return ParseWithLocale(input, DefaultLocale)
}

func ParseWithLocale(input string, locale Locale) (*Document, error) {
	parser := parser{
		lines:  splitLines(input),
		locale: locale,
	}
	return parser.parse()
}

type parser struct {
	lines  []string
	locale Locale
}

func (p *parser) parse() (*Document, error) {
	doc := &Document{}
	var currentField *Field

	for i := 0; i < len(p.lines); i++ {
		lineNo := i + 1
		raw := p.lines[i]
		trimmedLeft := strings.TrimLeft(raw, " \t\r")
		if trimmedLeft == "" {
			continue
		}
		if strings.HasPrefix(trimmedLeft, ">") {
			continue
		}
		if strings.HasPrefix(trimmedLeft, "--") {
			embed, next, err := p.parseEmbed(i)
			if err != nil {
				return nil, err
			}
			doc.Elements = append(doc.Elements, embed)
			currentField = nil
			i = next
			continue
		}
		if strings.HasPrefix(trimmedLeft, "-") {
			if strings.HasPrefix(trimmedLeft, "--") {
				continue
			}
			item, err := parseItem(raw, lineNo)
			if err != nil {
				return nil, err
			}
			if currentField == nil {
				return nil, newError(lineNo, p.locale.ItemOutsideField(lineNo, raw))
			}
			if currentField.content == fieldContentNone {
				currentField.content = fieldContentItems
			}
			if currentField.content != fieldContentItems {
				return nil, newError(currentField.lineNumber, p.locale.MixedFieldContent(currentField.lineNumber))
			}
			currentField.items = append(currentField.items, item)
			continue
		}
		if strings.HasPrefix(trimmedLeft, "`") {
			element, kind, err := p.parseEscaped(raw, lineNo, currentField)
			if err != nil {
				return nil, err
			}
			if kind == "root" {
				doc.Elements = append(doc.Elements, element.(Root))
				if field, ok := element.(*Field); ok {
					currentField = field
				} else {
					currentField = nil
				}
			}
			continue
		}
		if strings.HasPrefix(trimmedLeft, ":") {
			return nil, newError(lineNo, p.locale.FieldWithoutKey(lineNo, raw))
		}
		if strings.HasPrefix(trimmedLeft, "=") {
			return nil, newError(lineNo, p.locale.AttributeWithoutKey(lineNo, raw))
		}

		element, kind, err := p.parsePlain(raw, lineNo, currentField)
		if err != nil {
			return nil, err
		}
		if kind == "root" {
			doc.Elements = append(doc.Elements, element.(Root))
			if field, ok := element.(*Field); ok {
				currentField = field
			} else {
				currentField = nil
			}
		}
	}

	return doc, nil
}

func (p *parser) parseEmbed(start int) (*Embed, int, error) {
	lineNo := start + 1
	raw := p.lines[start]
	trimmedLeft := strings.TrimLeft(raw, " \t\r")
	dashes := countPrefix(trimmedLeft, '-')
	rest := strings.TrimSpace(trimmedLeft[dashes:])
	if rest == "" {
		return nil, 0, newError(lineNo, p.locale.EmbedWithoutKey(lineNo, raw))
	}
	valueLines := make([]string, 0)
	for i := start + 1; i < len(p.lines); i++ {
		candidate := strings.TrimLeft(p.lines[i], " \t\r")
		candidateDashes := countPrefix(candidate, '-')
		if candidateDashes == dashes && candidateDashes > 0 {
			candidateKey := strings.TrimSpace(candidate[candidateDashes:])
			if candidateKey == rest {
				var value *string
				if len(valueLines) > 0 {
					joined := strings.Join(valueLines, "\n")
					value = &joined
				}
				return &Embed{key: rest, lineNumber: lineNo, terminatorLine: i + 1, value: value}, i, nil
			}
		}
		valueLines = append(valueLines, p.lines[i])
	}
	return nil, 0, newError(lineNo, p.locale.UnterminatedEmbed(rest, lineNo))
}

func (p *parser) parseEscaped(raw string, lineNo int, currentField *Field) (any, string, error) {
	trimmedLeft := strings.TrimLeft(raw, " \t\r")
	ticks := countPrefix(trimmedLeft, '`')
	rest := trimmedLeft[ticks:]
	rest = strings.TrimLeft(rest, " \t\r")
	key, tail, ok := readEscapedKey(rest, ticks)
	if !ok {
		return nil, "", newError(lineNo, p.locale.UnterminatedEscapedKey(lineNo, raw))
	}
	if strings.TrimSpace(key) == "" {
		return nil, "", newError(lineNo, p.locale.EscapeWithoutKey(lineNo, raw))
	}
	return p.parseAfterKey(strings.TrimSpace(key), tail, raw, lineNo, currentField)
}

func (p *parser) parsePlain(raw string, lineNo int, currentField *Field) (any, string, error) {
	trimmedLeft := strings.TrimLeft(raw, " \t\r")
	operatorIdx := strings.IndexAny(trimmedLeft, ":=")
	if operatorIdx == -1 {
		key := strings.TrimRight(trimmedLeft, " \t\r")
		return &Flag{key: key, lineNumber: lineNo}, "root", nil
	}
	key := strings.TrimRight(trimmedLeft[:operatorIdx], " \t\r")
	tail := trimmedLeft[operatorIdx:]
	return p.parseAfterKey(key, tail, raw, lineNo, currentField)
}

func (p *parser) parseAfterKey(key, tail, raw string, lineNo int, currentField *Field) (any, string, error) {
	if key == "" {
		if len(tail) > 0 && tail[0] == '=' {
			return nil, "", newError(lineNo, p.locale.AttributeWithoutKey(lineNo, raw))
		}
		return nil, "", newError(lineNo, p.locale.FieldWithoutKey(lineNo, raw))
	}
	trimmedTail := strings.TrimLeft(tail, " \t\r")
	if trimmedTail == "" {
		return &Flag{key: key, lineNumber: lineNo}, "root", nil
	}
	switch trimmedTail[0] {
	case ':':
		value := trimToken(trimmedTail[1:])
		field := &Field{key: key, lineNumber: lineNo}
		if value != nil {
			field.content = fieldContentValue
			field.value = value
		}
		return field, "root", nil
	case '=':
		if currentField == nil {
			return nil, "", newError(lineNo, p.locale.AttributeOutsideField(lineNo, raw))
		}
		if currentField.content == fieldContentNone {
			currentField.content = fieldContentAttributes
		}
		if currentField.content != fieldContentAttributes {
			return nil, "", newError(currentField.lineNumber, p.locale.MixedFieldContent(currentField.lineNumber))
		}
		attribute := &Attribute{key: key, lineNumber: lineNo, value: trimToken(trimmedTail[1:])}
		currentField.attributes = append(currentField.attributes, attribute)
		return attribute, "nested", nil
	default:
		return nil, "", newError(lineNo, p.locale.InvalidAfterEscape(lineNo, raw))
	}
}

func parseItem(raw string, lineNo int) (*Item, error) {
	trimmedLeft := strings.TrimLeft(raw, " \t\r")
	return &Item{lineNumber: lineNo, value: trimToken(trimmedLeft[1:])}, nil
}

func splitLines(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, "\n")
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func countPrefix(s string, ch byte) int {
	count := 0
	for count < len(s) && s[count] == ch {
		count++
	}
	return count
}

func trimToken(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func readEscapedKey(rest string, ticks int) (string, string, bool) {
	var key strings.Builder
	for i := 0; i < len(rest); {
		if rest[i] == '`' {
			count := countPrefix(rest[i:], '`')
			if count == ticks {
				return key.String(), rest[i+ticks:], true
			}
		}
		key.WriteByte(rest[i])
		i++
	}
	return "", "", false
}
