package eno

import "testing"

func TestParseFieldValue(t *testing.T) {
	doc, err := Parse("field: value")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	field, err := doc.Field("field")
	if err != nil {
		t.Fatalf("field lookup failed: %v", err)
	}
	value, err := field.RequiredValueString()
	if err != nil {
		t.Fatalf("value lookup failed: %v", err)
	}
	if value != "value" {
		t.Fatalf("got %q", value)
	}
}

func TestParseFieldAttributes(t *testing.T) {
	doc, err := Parse("field:\nattribute = value")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	field, _ := doc.Field("field")
	attributes, err := field.Attributes()
	if err != nil {
		t.Fatalf("attributes failed: %v", err)
	}
	if len(attributes) != 1 {
		t.Fatalf("got %d attributes", len(attributes))
	}
	if attributes[0].Key() != "attribute" {
		t.Fatalf("got key %q", attributes[0].Key())
	}
	value, err := attributes[0].RequiredValueString()
	if err != nil || value != "value" {
		t.Fatalf("got value %q err %v", value, err)
	}
}

func TestParseFieldItems(t *testing.T) {
	doc, err := Parse("field:\n- item1\n- item2")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	field, _ := doc.Field("field")
	items, err := field.Items()
	if err != nil {
		t.Fatalf("items failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items", len(items))
	}
	value, err := items[1].RequiredValueString()
	if err != nil || value != "item2" {
		t.Fatalf("got value %q err %v", value, err)
	}
}

func TestParseEmbed(t *testing.T) {
	doc, err := Parse("-- embed\nvalue\n-- embed")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	embed, err := doc.Embed("embed")
	if err != nil {
		t.Fatalf("embed lookup failed: %v", err)
	}
	value, err := embed.RequiredValueString()
	if err != nil || value != "value" {
		t.Fatalf("got value %q err %v", value, err)
	}
}

func TestParseEscapedFlag(t *testing.T) {
	doc, err := Parse("`` `flag` ``")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	flag, err := doc.Flag("`flag`")
	if err != nil {
		t.Fatalf("flag lookup failed: %v", err)
	}
	if flag.LineNumber() != 1 {
		t.Fatalf("got line %d", flag.LineNumber())
	}
}

func TestFieldRequiredValueInt(t *testing.T) {
	doc, err := Parse("field: 23")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	field, _ := doc.Field("field")
	value, err := field.RequiredValueInt()
	if err != nil || value != 23 {
		t.Fatalf("got value %d err %v", value, err)
	}
}

func TestFieldRequiredValueIntError(t *testing.T) {
	doc, err := Parse("field: thirtytwo")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	field, _ := doc.Field("field")
	_, err = field.RequiredValueInt()
	if err == nil {
		t.Fatal("expected error")
	}
	parseErr := err.(*Error)
	if parseErr.Line != 1 || parseErr.Message != "strconv.Atoi: parsing \"thirtytwo\": invalid syntax" {
		t.Fatalf("got %#v", parseErr)
	}
}

func TestFieldAttributeRequiredValueString(t *testing.T) {
	doc, err := Parse("> comment\nfield:\nattribute = value")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	field, _ := doc.Field("field")
	attribute, err := field.Attribute("attribute")
	if err != nil {
		t.Fatalf("attribute failed: %v", err)
	}
	value, err := attribute.RequiredValueString()
	if err != nil || value != "value" {
		t.Fatalf("got value %q err %v", value, err)
	}
}

func TestFieldAttributeRequiredValueStringError(t *testing.T) {
	doc, err := Parse("> comment\nfield:\nattribute =")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	field, _ := doc.Field("field")
	attribute, _ := field.Attribute("attribute")
	_, err = attribute.RequiredValueString()
	if err == nil {
		t.Fatal("expected error")
	}
	parseErr := err.(*Error)
	if parseErr.Line != 3 || parseErr.Message != "Missing value" {
		t.Fatalf("got %#v", parseErr)
	}
}

func TestParsingErrors(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		line    int
		message string
	}{
		{"attribute outside field", "attribute = value", 1, "The attribute in line 1 is not contained within a field. ('attribute = value')"},
		{"attribute without key", "= value", 1, "The attribute in line 1 has no key. ('= value')"},
		{"embed without key", "--", 1, "The embed in line 1 has no key. ('--')"},
		{"escape without key", "` `", 1, "The escape sequence in line 1 specifies no key. ('` `')"},
		{"field without key", ": value", 1, "The field in line 1 has no key. (': value')"},
		{"invalid after escape", "`key` value", 1, "The escape sequence in line 1 can only be followed by an attribute or field operator. ('`key` value')"},
		{"item outside field", "- item", 1, "The item in line 1 is not contained within a field. ('- item')"},
		{"mixed field content", "field:\nattribute = value\n- item", 1, "The field in line 1 must contain either only attributes, only items, or only a value."},
		{"unterminated embed", "-- embed\n...", 1, "The embed 'embed' starting in line 1 is not terminated until the end of the document."},
		{"unterminated escaped key", "`key", 1, "The key escape sequence in line 1 is not terminated before the end of the line. ('`key')"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.input)
			if err == nil {
				t.Fatal("expected error")
			}
			parseErr := err.(*Error)
			if parseErr.Line != tc.line || parseErr.Message != tc.message {
				t.Fatalf("got %#v", parseErr)
			}
		})
	}
}

func TestPrettyPrint(t *testing.T) {
	doc, err := Parse("field:\nattribute = value\n\nflag\n\n-- embed\nline 1\nline 2\n-- embed")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	got := doc.PrettyPrint()
	want := "field:\nattribute = value\n\nflag\n\n-- embed\nline 1\nline 2\n-- embed"
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestPrettyPrintEscapesKeys(t *testing.T) {
	doc, err := Parse("` weird:key `: value\n\n`` `flag` ``")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	got := doc.PrettyPrint()
	want := "` weird:key `: value\n\n`` `flag` ``"
	if got != want {
		t.Fatalf("got %q", got)
	}
}
