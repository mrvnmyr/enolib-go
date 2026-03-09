# go-eno

`go-eno` is a small Go library and CLI for parsing [eno](https://eno-lang.org/) documents.

It currently supports:

- flags
- fields with values
- fields with attributes
- fields with items
- embeds
- escaped keys
- comment lines (ignored while parsing)

## Install

Build the CLI:

```sh
./build.sh
```

Run the test suite:

```sh
./tests.sh
```

## Example `.eno`

```eno
title: Example document

server:
host = localhost
port = 8080

features:
- prettyprint
- cli

enabled

-- description
This document is parsed by go-eno.
It can also be normalized back to eno text.
-- description
```

## Use As A Library

```go
package main

import (
	"fmt"
	"log"

	eno "go-eno"
)

func main() {
	input := `
title: Example document

server:
host = localhost
port = 8080

enabled

-- description
This document is parsed by go-eno.
-- description
`

	doc, err := eno.Parse(input)
	if err != nil {
		log.Fatal(err)
	}

	title, err := doc.Field("title")
	if err != nil {
		log.Fatal(err)
	}

	titleValue, err := title.RequiredValueString()
	if err != nil {
		log.Fatal(err)
	}

	server, err := doc.Field("server")
	if err != nil {
		log.Fatal(err)
	}

	host, err := server.Attribute("host")
	if err != nil {
		log.Fatal(err)
	}

	hostValue, err := host.RequiredValueString()
	if err != nil {
		log.Fatal(err)
	}

	enabled, err := doc.Flag("enabled")
	if err != nil {
		log.Fatal(err)
	}

	description, err := doc.Embed("description")
	if err != nil {
		log.Fatal(err)
	}

	descriptionValue, err := description.RequiredValueString()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("title:", titleValue)
	fmt.Println("host:", hostValue)
	fmt.Println("enabled at line:", enabled.LineNumber())
	fmt.Println("description:", descriptionValue)
	fmt.Println()
	fmt.Println(doc.PrettyPrint())
}
```

## Use The CLI

Read from a file:

```sh
./bin/eno config.eno
```

Read from stdin:

```sh
cat config.eno | ./bin/eno
```

The CLI parses the input and writes normalized eno back to stdout via `PrettyPrint()`.
