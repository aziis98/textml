# Text ML

My personal _text markup language_. Mostly use this as an alternative to Markdown as I prefer a more structured/extensible language for writing content (And I will never use/setup [MDX](https://mdxjs.com/)).

## Introduction

A **document** is composed of various _blocks_ that can also be nested. A **block** can be a of _text node_ or an _element node_. For example a document can be described as follows

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
            <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Debitis, odio.</p>

            #my.button-primary{ Ok }
            #my.button{ Other }
        </main>
    }
}
```
