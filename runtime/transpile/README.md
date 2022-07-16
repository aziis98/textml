# TextML / Transpile

## Usage

`textml transpile` is used for reading `*.tml` files and translating it to other common formats. For now you must pass a file and use one of the following command line options

-   `--format`, `-f`: Set a format, available formats are

    -   `go`: Uses <https://github.com/alecthomas/repr/> to show the parsed tree structure

    -   `json`: Converts the parsed document to JSON

    -   `inline-json`: As previous but inlined

    -   `transpile.html`: A simple semantic to convert `#html.ELEMENT { ... }` to the corresponding HTML element. This will get a major write pretty soon.

-   `--output`, `-o`: Set output file or "`-`" for stdout.

## Formats

### Document

Document is a Markdown like format that transpiles TextML to HTML.

#### Directives

- `#metadata` take a dictionary of TextML key values represented as a list of `#KEY { VALUE }` entries. For now this format is under-specified, available values are

    - _strings_: normal TextML text nodes.
    
    - _dict_: `#dict{ #KEY_1 { VALUE-1 } ... #KEY-N { VALUE-N } }`

- **Headings.**

    - `#title{ TEXT }`: `<h1>`
    - `#subtitle{ TEXT }`: `<h2>`
    - `#subsubtitle{ TEXT }`: `<h3>`
    - `#subsubsubtitle{ TEXT }`: `<h4>`

- **Formatting.**

    Common:

    - `#bold{ TEXT }`
    - `#italic{ TEXT }`
    - `#underline{ TEXT }`
    - `#strikethrough{ TEXT }`
    - `#code{ TEXT }`
    
    Hyperlinking and interactivity:
    
    - `#link{ TEXT }{ URL }`: the mnemonic is "_from_ TEXT _to_ URL".
    
    - `#note{ TEXT }{ NOTE_TEXT }`: like a footnote but more generic, wraps TEXT in a span with a random UUID. NOTE_TEXT is put in a `<div>` at the end of the rendered document with a reference to the generated UUID (something like `data-ref-note-id="..."`)
    
    - `#ref{ TEXT }{ REF }`: like previous but only generates the span with a data attribute like `data-ref="REF"` for interlinking. (This can be used for example for highlighting (with js) all spans with the same REF when one of them is hovered)

    Other:
    
    - `#color{ CSS_COLOR }{ TEXT }`: this utility wraps text in a span with the specified color.
    - `#mark{ CLASS_NAME }{ TEXT }`: this utility wraps text in a span with the specified css class.

- **Ordered and Unordered Lists.**

    - `#itemize{ #item{ TEXT } #item{ TEXT } ... }`
    - `#enumerate{ #item{ TEXT } #item{ TEXT } ... }`

- **Images.**

    - `#image{}`