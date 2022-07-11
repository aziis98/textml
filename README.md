# Text ML

My personal _text markup language_.

## Introduction

A **document** is composed of various nested _blocks_. A **block** can be a piece of _text block_ or a _node block_. For example

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

    This is the first block. Lorem ipsum dolor sit amet #bold { consectetur } adipisicing elit. Vitae architecto commodi officia natus ipsam labore fugit nisi quis. Deserunt, dolor #italic { consectetur nisi } placeat repellat velit #strikethrough { assumenda vero sed tenetur } hic.

    #subtitle { This is some title }

    This is the second block. Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae architecto commodi officia natus ipsam labore fugit nisi quis. Deserunt, #link { dolor consectetur nisi }{ https://en.wikipedia.org/wiki/Main_Page } placeat repellat velit assumenda vero sed tenetur hic.

    #code {
        #format { js }

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
