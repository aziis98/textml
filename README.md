# Text ML

My personal _text markup language_.

## Introduction

A **document** is composed of various nested _blocks_. A **block** can be a piece of _text block_ or a _node block_. For example

```
#document {
    #title { This is a short title }

    This is a paragraph with some text and #bold{ bold } and #italic{ italic } formatting.

    Elements can also have multiple argument as #link{this link to wikipedia}{https://en.wikipedia.org/}

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

## As a templating language

Let's say we want to use this as a templating language for HTML.

```
#layout{ base }{
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8" />
            <meta http-equiv="X-UA-Compatible" content="IE=edge" />
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <title>Example Page</title>

            #slot { head }
        </head>
        <body>
            #slot { body }
        </body>
    </html>
}

#page{ index.html }{
    #layout { base }

    #head {
        <link rel="stylesheet" href="styles/main.css">
    }

    #body {
        <main>
            <h1>Title</h1>
            <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Debitis, odio.</p>
        </main>
    }
}
```

## Composable Lexers

```go
type Lexer interface {
    Next() (rune, error)
    Peek() (rune, error)
    Done() bool
    Move(pos int) bool

    Expect(s string) error
}

func newPythonLikeLexer(tokens chan<- Token) {
    return indented.New(&indented.Config{
        // called on the "line" before the indentation increment
        PreIdentationStart: func(l Lexer) error {
            if err := l.Expect("block:"); err != nil {
                return err
            }

            tokens <- Token{ Type: "keyword", Value: "block:" }
            tokens <- Token{ Type: "block-begin", Value: "" }
            return nil
        },
        // called after the indentaion block ends (might be called multiple times if the indentation decreases by more than one level)
        IndentationEnd: func(l Lexer) error {
            tokens <- Token{ Type: "block-end", Value: "" }
        },
    })
}
```

<!--
### Python-variant

```
@document:
    @title: This is some title

    This is the first block. Lorem ipsum dolor sit amet @bold[consectetur] adipisicing elit. Vitae architecto commodi officia natus ipsam labore fugit nisi quis. Deserunt, dolor @italic[consectetur nisi] placeat repellat velit @strikethrough[assumenda vero sed tenetur] hic.

    @subtitle: This is some title

    This is the second block. Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae architecto commodi officia natus ipsam labore fugit nisi quis. Deserunt, @link[dolor consectetur nisi][https://en.wikipedia.org/wiki/Main_Page] placeat repellat velit assumenda vero sed tenetur hic.

    @code:
        @format: js
        function main() {
            console.log('Hello, World!')
        }

        main()

    To embed itself there is the multi-brace construct

    #code::
        #format[[ text-ml ]]

        This is some raw #code[text-ml]

    @list:
        @item: Item 1
        @item: Item 2
        @item: Item 3
        @item: Item 4
        @item:
            Item 5
        @item:
            Item 6
        @item:
            Item 7
        @item:
            Item 8
```

### C-variant

```
#document {
    #title { This is some title }

    This is the first block. Lorem ipsum dolor sit amet #bold { consectetur } adipisicing elit. Vitae architecto commodi officia natus ipsam labore fugit nisi quis. Deserunt, dolor #italic { consectetur nisi } placeat repellat velit #strikethrough{ assumenda vero sed tenetur } hic.

    #subtitle { This is some title }

    This is the second block. Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae architecto commodi officia natus ipsam labore fugit nisi quis. Deserunt, #link{ {dolor consectetur nisi} }{ https://en.wikipedia.org/wiki/Main_Page } placeat repellat velit assumenda vero sed tenetur hic.

    #code{
        #format { js }
        #theme { solaraized }

        function main() {
            console.log('Hello, World!')
        }

        main()
    }

    To embed itself there is the multi-brace construct

    #code {{
        #format {{ text-ml }}

        This is some raw #code { text-ml }
    }}

    -   Item 1
    -   Item 2
    -   Item 3

    #list {
        #item { Item 1 }
        #item { Item 2 }
        #item { Item 3 }
        #item { Item 4 }
        #item {
            Item 5
        }
        #item {
            Item 6
        }
        #item {
            Item 7
        }
        #item {
            Item 8
        }
    }

    #list-item {
        fiehffweoi
    }
    #list-item {
        fiehffweoi
    }
    #list-item {
        fiehffweoi

        #list-item {
            fiehffweoi
        }
        #list-item {
            fiehffweoi
        }
        #list-item {
            fiehffweoi
        }
    }
}
```

## Lisp-embedding

Clearly this is equivalent to _M-expressions_, for example

```
#defun{ sum }{ #params{ list } }{
    #if-else{ #eq{ #length{ list } }{ 0 } }{
        #of{ list }{ 0 }
    }{
        #plus{ #of{ list }{ 0 } }{ #sum{ #slice{ list }{ 1 } } }
    }
}

#set{ list-1 }{ #list{ 1 }{ 2 }{ 3 } }
#println{ #sum{ list-1 } }
```
-->
