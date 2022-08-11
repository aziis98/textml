# Textml / Template

This runtime provides a sub-language to define templates and layouts and combine them together. For example a basic layout file could be

_./example-layouts.tml_

```html
#template{ base-layout }{
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>#{ title }</title>
        
        <link rel="stylesheet" href="/assets/styles.css">

        #{ head }
    </head>
    <body>
        #{ body }
    </body>
    </html>
}

#template{ article-layout }{
    #extends{ base-layout }{
        #define{ title }{ Article - #{ article.title } }
        #define{ body }{
            <h1>#{ article.title }</h1>
            #{ article.body }
        }
    }
}
```

Then we can add a concrete page that uses this layout file

_./article-1.tml_

```html
#import{ ./example-layouts.tml }
#extends{ base-layout }{
    #define{ head }{
        <!-- maybe add KaTeX support just for this page -->
    }
    #define{ article.title }{ Article 1 }
    #define{ article.body }{
        <p>Some article</p>
    }
}

```

## Usage

`textml template [-o OUTPUT] FILES...` will evaluate each file in sequence starting from an empty context, the default `LoaderFunc` is `FileLoader` so `#import{ FILE }` will read and evaluate that file. The option `--output` or `-o` can be used to change the file to write to, by default its the `-` meaning stdout.

## Documentation

- `#template{ NAME }{ TEMPLATE }` defines a new template `NAME` with value `TEMPLATE`. Templates can be expanded using the `#extends{ NAME }{ ... }` directive

- `#define{ NAME }{ VALUE }` evaluates ast for `VALUE` when this define gets called and binds it to the variable `NAME`.

- `#{ EXPR }` evaluates the code inside or variable interpolation.

- `#extends{ NAME }{ BLOCK }` evaluates first the given `BLOCK` and then the template bind to `NAME`.

- `#import{ MODULE }` is used to "include" a module using the provider `LoaderFunc`, for example `FileLoader` reads a file and evaluates it in-place in the current context, the produced string is .

### Expressions

- `#char{ CHAR_NAME }`

- `#foreach{ ITEM }{ ITEMS }{ BLOCK }`

- `#intersperse{ ITEMS }{ SEPARATOR }`

- `#if{ CONDITION }{ IF_TRUE }`

- `#if{ CONDITION }{ IF_TRUE }{ IF_FALSE }`





