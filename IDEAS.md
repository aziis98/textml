# Ideas

List of some random ideas for using this language (top are more recent)

## Offload execution to other languages

For example the _document_ runtime could have a directive like `#import-js-preprocessor{ SCRIPT }` that processes an unknown definition in the document.

```js
import { transformInput, isElementNode, getTextContent } from 'textml'

/*
    #figure{ 
        #src{ IMAGE_PATH } 
        #placement{ POSITION }
        #description{ DESCRIPTION } 
    }
*/
transformInput(element => {
    const options = {}

    if (element.name === "figure") {
        element.args[0].filter(isElementNode).forEach(option => {
            options[option.name] = getTextContent(option.args[0])
        })
    }

    return `
        <div class="${['figure', option.placement].filter(Boolean).join(' ')}">
            <img src="${option.src}">
            <div class="description">
                ${option.description}
            </div>
        </div>   
    `
})
```

## Templating language (alternative)

One of the next thing I will start working on is a way to use this language as a templating language for building HTML pages.

```
#layout{ example-layout }{
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
    #use { example-layout }

    #define{ my.button }{
        <button class="button">#slot{}</button>
    }

    #define{ my.button-primary }{
        <button class="button primary">#slot{}</button>
    }

    #head {
        <link rel="stylesheet" href="styles/main.css">
    }

    #body {
        <main>
            <h1>Title</h1>
            <p>
                Lorem ipsum dolor sit amet consectetur adipisicing
                elit. Debitis, odio.
            </p>

            #my.button-primary{ Ok }
            #my.button{ Other }
        </main>
    }
}
```

## Multiple Runtimes

The CLI could support multiple runtimes like git with `textml RUNTIME`. Some examples are

-   `textml build`

    This runtime/environment reads a file called `build.tml` in the current directory and runs it like a `Makefile` (but somehow also with support for fan-out and fan-in).

    ```
    #rule{
        #inputs{ articles/*.tml }

        #input { articles/:name.tml }
        #output { dist/articles/:name.html }

        textml transpile -f html "#{ INPUT }" -o "#{ OUTPUT }"
    }

    #rule {
        #inputs { articles/*.tml }
        #output { dist/tags.json }

        ...extract from each article metadata its tags and create a "tags.json" file mapping each tag to a list of article ids with that tag...
    }

    #rule {
        #input { dist/tags.json }
        #outputs { dist/tags/*.html }

        ...for each tag in "tags.json" generate a corresponding page with links to every article containing that tag.
    }
    ```

-   `textml transpile -f FORMAT FILES...` to transpile files to other formats.

-   `textml query -e EXPRESSION FILES...` to query TML files with various selectors.

## Content Embedding / File Processing

```
#asset{ styles/main.css }{{
    html, body {
        margin: 0;
        min-height: 100vh;
    }
}}

#file{ example.txt }{
    #build.pipe-content{ fold -w 40 }

    Lorem ipsum dolor sit amet consectetur, adipisicing elit. Consequatur repellat cupiditate dolor ut, fuga recusandae quos vitae, ea necessitatibus dolores quam tenetur minima tempora eum numquam inventore suscipit vero consequuntur?
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

## As a tree/text transformation utility

```
#transform {
    #query{ #document > (#title: $title) }{
        <h1>#{ $title }</h1>
    }
    #query{ #document (( #item (+ #item)* ): $items) }{
        <ul>
        #foreach{ $item }{ $items }{
            <li>#{ $item }</li>
        }
        </ul>
    }
}
```

For the first query

```js
find($root, '#document', document => {
    findChild($document, '#title', title => {
        const $title = title.children
        replaceContent(title, ['<h1>', ...$title, '</h1>'])
    })
})
```

And something like this for the second one

```js
const onItems = $items => {
    replaceSequence($items, [
        '<ul>',
        ...$items.flatMap($item => {
            return ['<li>', ...$item.children, '</li>']
        })
        '</ul>',
    ])
}

find($root, '#document', document => {
    findNested($document, '#item', item => {
        const $items = []

        const findNextFn = item => {
            $items.push(item)
            const found = findAfter(item, '#item', findNextFn)

            if (!found) {
                onItems($items)
            }
        }

        findNextFn(item)
    })
})
```

## Composable Lexers

A way to make a composable lexer for the Python-like version

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
        PreIndentationStart: func(l Lexer) error {
            if err := l.Expect("block:"); err != nil {
                return err
            }

            tokens <- Token{ Type: "keyword", Value: "block:" }
            tokens <- Token{ Type: "block-begin", Value: "" }
            return nil
        },
        // called after the indentation block ends (might be called multiple times if the indentation decreases by more than one level)
        IndentationEnd: func(l Lexer) error {
            tokens <- Token{ Type: "block-end", Value: "" }
        },
    })
}
```

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
