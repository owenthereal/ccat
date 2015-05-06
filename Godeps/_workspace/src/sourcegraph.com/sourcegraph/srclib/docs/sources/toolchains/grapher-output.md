page_title: Grapher Output

# Grapher Output

Src will invoke the grapher, providing a JSON representation of a source unit (`*unit.SourceUnit`)
in through stdin.

## Output Schema

The output is a single JSON object with three fields that represent lists of
Definitions, References, and Documentation data respectively. This should be printed to stdout.

[[.code "graph/output.pb.go" "Output"]]

### Def Object Structure
[[.code "graph/def.pb.go" "Def "]]

### Ref Object Structure
[[.code "graph/ref.pb.go" "Ref"]]

### Docs Object Structure
[[.code "graph/doc.pb.go" "Doc"]]

## Example: Grapher output on [jashkenas/underscore](https://github.com/jashkenas/underscore)
```json
{
  "Defs": [
    {
      "Path": "commonjs/test/arrays.js",
      "TreePath": "-commonjs/test/arrays.js",
      "Kind": "module",
      "Exported": true,
      "Data": {
        "Kind": "commonjs-module",
        "Key": {
          "namespace": "commonjs",
          "module": "test/arrays.js",
          "path": ""
        },
        "jsgSymbolData": {
          "nodejs": {
            "moduleExports": true
          }
        },
        "Type": "{}",
        "IsFunc": false
      },
      "Name": "test/arrays",
      "File": "test/arrays.js",
      "DefStart": 0,
      "DefEnd": 0
    },
    ...
  ],
  "Refs" : [
    {
      "DefRepo": "",
      "DefUnitType": "",
      "DefUnit": "",
      "DefPath": "commonjs/underscore.js/-/union",
      "File": "test/arrays.js",
      "Start": 7610,
      "End": 7615
    },
    ...
  ],
  "Docs" : [
    {
      "Path": "commonjs/test/vendor/qunit.js/-/jsDump/parsers/functionArgs",
      "Format": "",
      "Data": "function calls it internally, it's the arguments part of the function"
    },
    ...
  ]
```
