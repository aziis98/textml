# TextML / Transpile

## Usage

`textml transpile` is used for reading `*.tml` files and translating it to other common formats. For now you must pass a file and use one of the following command line options

-   `--format`, `-f`: Set a format, available formats are

    -   `go`: Uses <https://github.com/alecthomas/repr/> to show the parsed tree structure

    -   `json`: Converts the parsed document to JSON

    -   `inline-json`: As previous but inlined

    -   `transpile.html`: A simple semantic to convert `#html.ELEMENT { ... }` to the corresponding HTML element. This will get a major write pretty soon.

-   `--output`, `-o`: Set output file or "`-`" for stdout.