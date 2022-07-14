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
	transpile	Used to read .tml files and convert them to other formats
	template	Use textml as a templating language
`

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString(usage)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "transpile":
		transpileCmd := flag.NewFlagSet("transpile", flag.ExitOnError)
		listFormats := transpileCmd.Bool("list-formats", false, "Display available formats")
		format := transpileCmd.StringP(
			"format", "f", "repr", `output format of the parsed file`,
		)
		output := transpileCmd.StringP(
			"output", "o", "-", `output file, "-" is stdout`,
		)
		transpileCmd.Parse(os.Args[2:])

		if *listFormats {
			fmt.Printf("Available formats:\n")
			for format, _ := range transpile.Registry {
				fmt.Printf("- %q\n", format)
			}

			os.Exit(0)
		}

		outputFile := os.Stdout
		if *output != "-" {
			f, err := os.Create(*output)
			if err != nil {
				log.Fatal(err)
			}

			outputFile = f
		}

		if transpileCmd.NArg() == 0 {
			log.Fatal("invalid number of arguments")
		}

		inputFile, err := os.Open(transpileCmd.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		commandTranspile(inputFile, outputFile, *format)
	case "template":
		templateCmd := flag.NewFlagSet("template", flag.ExitOnError)
		output := templateCmd.StringP(
			"output", "o", "-", `output file, "-" is stdout`,
		)
		templateCmd.Parse(os.Args[2:])

		outputFile := os.Stdout
		if *output != "-" {
			f, err := os.Create(*output)
			if err != nil {
				log.Fatal(err)
			}

			outputFile = f
		}

		inputFile, err := os.Open(templateCmd.Arg(0))
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
	if err := transpiler.Transpile(outputFile, doc); err != nil {
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
