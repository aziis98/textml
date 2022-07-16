package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aziis98/textml"
	"github.com/aziis98/textml/runtime/template"
	"github.com/aziis98/textml/runtime/transpile"

	flag "github.com/spf13/pflag"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lshortfile | log.Ltime)
}

const usage = `usage: textml COMMAND ...

Commands:
	transpile   Used to read .tml files and convert them to other formats
	template    Use textml as a templating language
`

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString(usage)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "transpile":
		cmd := flag.NewFlagSet("transpile", flag.ExitOnError)
		cmd.Usage = func() {
			fmt.Printf("usage: textml transpile [-f FORMAT] FILE\n\n")
			cmd.PrintDefaults()
		}

		// Options
		var listFormats bool
		cmd.BoolVar(&listFormats, "list-formats", false, "Display available formats")

		var format string
		cmd.StringVarP(&format, "format", "f", "repr", `output format of the parsed file`)

		var output string
		cmd.StringVarP(&output, "output", "o", "-", `output file, "-" is stdout`)

		showHelp := false
		cmd.BoolVarP(&showHelp, "help", "h", false, "Display help text")

		if err := cmd.Parse(os.Args[2:]); err != nil {
			if err != flag.ErrHelp {
				log.Fatal(err)
			}
		}

		if showHelp || cmd.NArg() == 0 {
			cmd.Usage()
			os.Exit(0)
		}

		if listFormats {
			fmt.Printf("Available formats:\n")
			for format := range transpile.Registry {
				fmt.Printf("- %q\n", format)
			}

			os.Exit(0)
		}

		outputFile := os.Stdout
		if output != "-" {
			f, err := os.Create(output)
			if err != nil {
				log.Fatal(err)
			}

			outputFile = f
		}

		inputFile, err := os.Open(cmd.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		commandTranspile(inputFile, outputFile, format)
	case "template":
		cmd := flag.NewFlagSet("template", flag.ExitOnError)
		cmd.Usage = func() {
			fmt.Printf("usage: textml template [--output|-o OUTPUT] FILES...\n\n")
			cmd.PrintDefaults()
		}

		var output string
		cmd.StringVarP(&output, "output", "o", "-", `output file, "-" is stdout`)

		var showHelp bool
		cmd.BoolVarP(&showHelp, "help", "h", false, "Display help text")

		if err := cmd.Parse(os.Args[2:]); err != nil {
			if err != flag.ErrHelp {
				log.Fatal(err)
			}
		}

		if showHelp || cmd.NArg() == 0 {
			cmd.Usage()
			os.Exit(0)
		}

		outputFile := os.Stdout
		if output != "-" {
			f, err := os.Create(output)
			if err != nil {
				log.Fatal(err)
			}

			outputFile = f
		}

		inputFile, err := os.Open(cmd.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		commandTemplate(inputFile, outputFile)
	default:
		log.Fatalf("invalid command %q", os.Args[1])
	}
}

func commandTranspile(inputFile *os.File, outputFile *os.File, format string) {
	doc, err := textml.ParseDocument(bufio.NewReader(inputFile))
	if err != nil {
		log.Fatal(err)
	}

	transpiler := transpile.Registry[format]
	if err := transpiler.Transpile(doc, outputFile); err != nil {
		log.Fatal(err)
	}
}

func commandTemplate(inputFile *os.File, outputFile *os.File) {
	doc, err := textml.ParseDocument(bufio.NewReader(inputFile))
	if err != nil {
		log.Fatal(err)
	}

	ctx := template.New(&template.Config{
		LoaderFunc: template.FileLoader,
	})

	s, err := ctx.Evaluate(doc)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := outputFile.WriteString(
		strings.TrimSpace(s) + "\n",
	); err != nil {
		log.Fatal(err)
	}
}
