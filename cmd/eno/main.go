package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	eno "go-eno"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("eno", flag.ContinueOnError)
	flags.SetOutput(stderr)

	if err := flags.Parse(args); err != nil {
		return 2
	}

	input, err := readInput(flags.Args(), stdin)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	document, err := eno.Parse(input)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if _, err := io.WriteString(stdout, document.PrettyPrint()); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	return 0
}

func readInput(args []string, stdin io.Reader) (string, error) {
	switch len(args) {
	case 0:
		data, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		return string(data), nil
	case 1:
		if args[0] == "-" {
			data, err := io.ReadAll(stdin)
			if err != nil {
				return "", fmt.Errorf("read stdin: %w", err)
			}
			return string(data), nil
		}
		data, err := os.ReadFile(args[0])
		if err != nil {
			return "", fmt.Errorf("read %s: %w", args[0], err)
		}
		return string(data), nil
	default:
		return "", fmt.Errorf("usage: eno [file|-]")
	}
}
