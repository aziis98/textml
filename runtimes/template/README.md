# Textml / Template

This runtime provides a sub-language to define templates and layouts and combine them together. For example a basic layout file could be

_./example-layouts.tml_

```html
#define{ base-layout }{
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

#define{ article-layout }{
    #define{ title }{ Article - #{ article.title } }
    #define{ body }{
        <h1>#{ article.title }</h1>
        #{ article.body }
    }
    
    #include{ base-layout }
}
```

Then we can add a concrete page that uses this layout file

_./article-1.tml_

```html
#import{ ./example-layouts.tml }

#output{
    #define{ head }{
        <!-- maybe add KaTeX support just for this page -->
    }

    #define{ article.title }{ Article 1 }
    #define{ article.body }{
        <p>Some article</p>
    }

    #include{ base-layout }
}
```

## Documentation

- `#layout{ NAME }{  }`