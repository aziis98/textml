package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/alecthomas/repr"
	"github.com/aziis98/go-text-ml/lexer"
	"github.com/aziis98/go-text-ml/parser"
	"github.com/aziis98/go-text-ml/transpiler"
	flag "github.com/spf13/pflag"
)

var (
	format string
	output string
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lshortfile | log.Ltime)
}

func main() {
	//
	// Flags
	//
	flag.StringVarP(&format, "format", "f", "go",
		`output format of the parsed file: go, json, inline-json, transpile.html`,
	)
	flag.StringVarP(&output, "output", "o", "-",
		`output file, "-" is stdout`,
	)
	flag.Parse()

	outWriter := os.Stdout
	if output != "-" {
		f, err := os.Create(output)
		if err != nil {
			log.Fatal(err)
		}

		outWriter = f
	}

	if flag.NArg() == 0 {
		log.Fatal("not enough arguments, must pass at least *.tml file to parse")
	}

	file := flag.Arg(0)

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	l := lexer.New(bufio.NewReader(f))

	tokens, err := l.AllTokens()
	if err != nil {
		log.Fatal(err)
	}

	doc, err := parser.ParseDocument(tokens)
	if err != nil {
		log.Fatal(err)
	}

	switch format {
	case "go":
		repr.New(outWriter).Println(doc)
	case "json":
		bytes, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := outWriter.Write(bytes); err != nil {
			log.Fatal(err)
		}
		if _, err := outWriter.WriteString("\n"); err != nil {
			log.Fatal(err)
		}
	case "inline-json":
		bytes, err := json.Marshal(doc)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := outWriter.Write(bytes); err != nil {
			log.Fatal(err)
		}
		if _, err := outWriter.WriteString("\n"); err != nil {
			log.Fatal(err)
		}
	case "transpile.html":
		log.Printf("Transpiling to HTML...")

		htmlTranspiler := &transpiler.Html{Inline: false}

		if err := htmlTranspiler.Transpile(outWriter, doc); err != nil {
			log.Fatal(err)
		}
	}
}
