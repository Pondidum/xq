# XQ
*Like [JQ](https://stedolan.github.io/jq/), but for XML*

## Methods

* `read` - Execute an XPath query against a document

### Read

Usage
```bash
xq read [options] <xpath> <file_path>
```

You can read from a file directly:
```bash
xq read 'count(//book[@type="short"])' books.xml
```

Or piped to stdin, by using `-` as the filepath:
```bash
cat books.xml | xq read 'count(//book[@type="short"])' -
```

