# textml

My personal textual markup language. Mostly a more structured alternative to Markdown as I prefer a more extensible language for writing content (and I will never use/setup [MDX](https://mdxjs.com/)).

## Introduction

A **document** is composed of various _blocks_ that can also be nested. A **block** can be a of _text node_ or an _element node_. For example a document can be described as follows

```
#document {
    #title { This is a short title }

    This is a paragraph with some text and #bold{ bold } and
    #italic{ italic } formatting.

    Elements can also have multiple argument as #link{this link
    to wikipedia}{https://en.wikipedia.org/}

    #subtitle { Another section }

    Code blocks can be easily nested and annotated, for example

    #code {{
        #format {{ textml }}

        #document {
            #title { Example }

            With #underline{some text}
        }
    }}

    As long as braces are balanced any meta-depth can be reached.

    #code {{{
        #code {{
            #code {
                let x = 1;
            }
        }}
    }}}
}
```

## Usage

For now there is a small CLI for working with the various "runtimes"

- `textml transpile OPTIONS...`

    Used to transpile TextML to other structured formats like json or html.

- `textml template [--output|-o OUTPUT] FILES...`

    Used to interpret TextML files as templates, for now the only supported directives are `#define{ NAME }{ TEMPLATE }`, `#{ NAME }`, `#import{ FILE }`, `#extends{ NAME }`.

