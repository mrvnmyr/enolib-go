package eno

import "strings"

func (d *Document) PrettyPrint() string {
	var builder strings.Builder

	for index, element := range d.Elements {
		if index > 0 {
			builder.WriteString("\n\n")
		}
		writeRoot(&builder, element)
	}

	return builder.String()
}

func writeRoot(builder *strings.Builder, root Root) {
	switch element := root.(type) {
	case *Field:
		writeField(builder, element)
	case *Flag:
		builder.WriteString(formatKey(element.key))
	case *Embed:
		key := formatKey(element.key)
		builder.WriteString("-- ")
		builder.WriteString(key)
		if element.value != nil {
			builder.WriteString("\n")
			builder.WriteString(*element.value)
		}
		builder.WriteString("\n-- ")
		builder.WriteString(key)
	}
}

func writeField(builder *strings.Builder, field *Field) {
	builder.WriteString(formatKey(field.key))
	builder.WriteString(":")
	if field.value != nil {
		builder.WriteString(" ")
		builder.WriteString(*field.value)
		return
	}
	for _, attribute := range field.attributes {
		builder.WriteString("\n")
		builder.WriteString(formatKey(attribute.key))
		builder.WriteString(" =")
		if attribute.value != nil {
			builder.WriteString(" ")
			builder.WriteString(*attribute.value)
		}
	}
	for _, item := range field.items {
		builder.WriteString("\n-")
		if item.value != nil {
			builder.WriteString(" ")
			builder.WriteString(*item.value)
		}
	}
}

func formatKey(key string) string {
	if !needsEscaping(key) {
		return key
	}

	ticks := 1
	for {
		delimiter := strings.Repeat("`", ticks)
		if !strings.Contains(key, delimiter) {
			return delimiter + " " + key + " " + delimiter
		}
		ticks++
	}
}

func needsEscaping(key string) bool {
	if key == "" {
		return true
	}
	if strings.TrimSpace(key) != key {
		return true
	}
	if strings.ContainsAny(key, "\n\r:=") {
		return true
	}
	switch key[0] {
	case '-', '>', '`', ':', '=':
		return true
	}
	return false
}
